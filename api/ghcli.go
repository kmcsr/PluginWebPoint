
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var githubClient = initGithubClient()

type GToken struct{
	Token string `json:"access_token"`
	Typ   string `json:"token_type"`
	Scope string `json:"scope"`
}

func (t GToken)GetAuth()(string){
	return "Bearer " + t.Token
}

type ghClient struct{
	cli *http.Client
	appId     string
	appSecret string
	token *GToken
}

func initGithubClient()(c *ghClient){
	c = &ghClient{
		cli: &http.Client{
			Timeout: time.Second * 5,
		},
		appId: os.Getenv("GH_CLI_ID"),
		appSecret: os.Getenv("GH_CLI_SEC"),
	}
	if len(c.appId) > 0 {
		if os.Getenv("GH_OAUTH") == "true" {
			loger.Infof("Pending github OAuth...")
			err := c.ghOAuth(func(code string, uri string)(error){
				loger.Warnf("Please use code '%s' at <%s>", code, uri)
				return nil
			})
			if err != nil {
				loger.Errorf("Cannot pass Github OAuth: %v", err)
			}else{
				loger.Infof("Github OAuth done!")
			}
		}else if len(c.appSecret) == 0 {
			loger.Warnf("Github auth skipped, set `GH_CLI_SEC` as your github application secret to auth to github.")
		}
	}else{
		loger.Warnf("Github auth skipped, set `GH_CLI_ID`, `GH_CLI_SEC` as your github application id and secret to enable it.")
	}
	c.pingGhApi()
	return
}

func (c *ghClient)pingGhApi()(err error){
	const entryPoint = "https://api.github.com/"

	var res *http.Response
	if res, err = c.Get(entryPoint); err != nil {
		return
	}
	res.Body.Close()
	loger.Infof("Status: %s", res.Status)
	loger.Infof("x-ratelimit-limit:     %s", res.Header.Get("x-ratelimit-limit"))
	loger.Infof("x-ratelimit-remaining: %s", res.Header.Get("x-ratelimit-remaining"))
	loger.Infof("x-ratelimit-used:      %s", res.Header.Get("x-ratelimit-used"))
	ts, _ := strconv.ParseInt(res.Header.Get("x-ratelimit-reset"), 10, 64)
	loger.Infof("x-ratelimit-reset: %s", time.Unix(ts, 0))
	return
}

func (c *ghClient)ghOAuth(callback func(code string, uri string)(error))(err error){
	const codePoint = "https://github.com/login/device/code"
	const pollPoint = "https://github.com/login/oauth/access_token"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute * 5)
	defer cancel()

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, "POST", codePoint, strings.NewReader("client_id=" + c.appId)); err != nil {
		return
	}
	req.Header.Set("User-Agent", "PluginWebPoint-App")
	req.Header.Set("Accept", "application/json")

	var res *http.Response
	if res, err = c.cli.Do(req); err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return &StatusCodeErr{ Code: res.StatusCode }
	}
	var codePayload struct{
		DeviceCode      string `json:"device_code"`
		UserCode        string `json:"user_code"`
		VerificationURI string `json:"verification_uri"`
		ExpiresIn       int    `json:"expires_in"`
		Interval        int    `json:"interval"`
	}
	var data []byte
	if data, err = io.ReadAll(res.Body); err != nil {
		return
	}
	if err = json.Unmarshal(data, &codePayload); err != nil {
		return
	}
	if err = callback(codePayload.UserCode, codePayload.VerificationURI); err != nil {
		return
	}
	if req, err = http.NewRequestWithContext(ctx, "POST", pollPoint,
		strings.NewReader(fmt.Sprintf("client_id=%s&device_code=%s&grant_type=%s",
			c.appId, codePayload.DeviceCode,
			"urn:ietf:params:oauth:grant-type:device_code"))); err != nil {
		return
	}

	req.Header.Set("User-Agent", "PluginWebPoint-App")
	req.Header.Set("Accept", "application/json")

	var pollPayload struct{
		GToken
		Err     string `json:"error"`
		ErrDesc string `json:"error_description"`
		ErrUri  string `json:"error_uri"`
	}
	interval := (time.Duration)(codePayload.Interval) * time.Second
	for {
		req.Body, _ = req.GetBody()
		if res, err = c.cli.Do(req); err != nil {
			return
		}
		if res.StatusCode != http.StatusOK {
			res.Body.Close()
			return &StatusCodeErr{ Code: res.StatusCode }
		}
		data, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return
		}
		pollPayload.Err = ""
		if err = json.Unmarshal(data, &pollPayload); err != nil {
			return
		}
		if len(pollPayload.Err) > 0 {
			switch pollPayload.Err {
			case "slow_down":
				interval += 5 * time.Second
				fallthrough
			case "authorization_pending":
				loger.Debug("Authorization pending...")
				time.Sleep(interval)
				continue
			}
			err = fmt.Errorf("Error when polling: [%s]: %s; see <%s>", pollPayload.Err, pollPayload.ErrDesc, pollPayload.ErrUri)
			return
		}
		c.token = &pollPayload.GToken
		break
	}
	return
}

func (c *ghClient)Get(url string)(*http.Response, error){
	return c.GetWithContext(context.Background(), url)
}

func (c *ghClient)GetWithContext(ctx context.Context, url string)(res *http.Response, err error){
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "PluginWebPoint-App")
	if c.token != nil {
		req.Header.Set("Authorization", c.token.GetAuth())
	}else if len(c.appSecret) > 0 {
		req.SetBasicAuth(c.appId, c.appSecret)
	}
	if res, err = c.cli.Do(req); err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		if res.StatusCode != http.StatusNoContent {
			err = &StatusCodeErr{ Code: res.StatusCode }
		}
	}
	return
}
