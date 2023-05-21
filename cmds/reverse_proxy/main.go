
package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/kmcsr/go-logger"
	"github.com/kmcsr/go-logger/logrus"
)

var loger logger.Logger = initLogger()

func initLogger()(loger logger.Logger){
	loger = logrus.Logger
	if os.Getenv("DEBUG") == "true" {
		loger.SetLevel(logger.TraceLevel)
	}else{
		loger.SetLevel(logger.InfoLevel)
		_, err := logger.OutputToFile(loger, "/var/log/pwp/rpx/latest.log", os.Stdout)
		if err != nil {
			panic(err)
		}
	}
	return
}

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`

	Web  string            `json:"web"`
	APIs map[string]string `json:"apis"`
}

var config *Config = loadConfig()

func loadConfig()(cfg *Config){
	const configPath = "/etc/pwp/rpx/config.json"
	var data []byte
	var err error
	if data, err = os.ReadFile(configPath); err != nil {
		loger.Fatalf("Cannot read config at %s: %v", configPath, err)
	}
	cfg = new(Config)
	if err = json.Unmarshal(data, cfg); err != nil {
		loger.Fatalf("Cannot parse config at %s: %v", configPath, err)
	}
	return
}

func ProxyOf(prefix string, target string)(http.HandlerFunc){
	targetUrl, err := url.Parse(target)
	if err != nil {
		loger.Fatalf("Cannot parse target url %q: %v", target, err)
	}
	return func(rw http.ResponseWriter, r *http.Request){
		path := r.URL.Path

		ctx := r.Context()
		req := r.Clone(ctx)
		req.Host = targetUrl.Host
		req.RequestURI = ""
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host
		if len(prefix) > 0 {
			req.URL.Path = strings.TrimLeft(path, prefix)
			if len(req.URL.Path) == 0 || req.URL.Path[0] != '/' {
				req.URL.Path = "/" + req.URL.Path
			}
		}
		if r.ContentLength == 0 {
			req.Body = nil
		}
		if req.Body != nil {
			defer req.Body.Close()
		}

		loger.Debugf("Requesting %q -> %q", path, req.URL.String())
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			if res != nil {
				res.Body.Close()
			}
			loger.Errorf("Error when requesting %s: %v", req.URL.String(), err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write(([]byte)(err.Error()))
			return
		}
		defer res.Body.Close()
		for k, v := range res.Header {
			rw.Header()[k] = v
		}
		rw.WriteHeader(res.StatusCode)
		io.Copy(rw, res.Body)
	}
}

func main(){
	for prefix, target := range config.APIs {
		prefix = "/" + prefix + "/"
		loger.Infof("Proxy %q to %q", prefix, target)
		http.HandleFunc(prefix, ProxyOf(prefix, target))
	}

	webProxy := ProxyOf("", config.Web)
	http.HandleFunc("/assets/", webProxy)
	http.HandleFunc("/", webProxy)

	server := &http.Server{
		Addr: net.JoinHostPort(config.Host, strconv.Itoa(config.Port)),
		Handler: http.DefaultServeMux,
		ReadTimeout: time.Second * 30,
		WriteTimeout: time.Second * 60,
	}

	done := make(chan struct{}, 0)
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	go func(){
		defer close(done)
		loger.Infof("Server start at %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			loger.Fatal(err)
		}
	}()

	select {
	case <-sigch:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
		server.Shutdown(ctx)
		cancel()
	case <-done:
	}
}
