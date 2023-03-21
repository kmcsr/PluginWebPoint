
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kataras/iris/v12"
	irisContext "github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kmcsr/PluginWebPoint/api"
	"github.com/kmcsr/PluginWebPoint/api/mysqlimpl"
)

const DevApiVerion = 0

type ErrResp struct{
	Status string `json:"status"`
	Name   string `json:"error"`
	Msg    string `json:"message"`
	Extra  string `json:"extra,omitempty"`
}

func NewErrResp(name string, err error)(*ErrResp){
	return &ErrResp{
		Status: "error",
		Name: name,
		Msg: err.Error(),
	}
}

var apiIns api.API = nil
var ghCli = api.InitGithubCli()

func main(){
	address := ""
	if len(os.Args) >= 2 {
		address = os.Args[1]
	}

	username := os.Getenv("DB_USER")
	passwd := os.Getenv("DB_PASSWD")
	dbaddress := os.Getenv("DB_ADDR")
	database := os.Getenv("DB_NAME")

	apiIns = mysqlimpl.NewMySqlAPI(username, passwd, dbaddress, database, ghCli)

	app := iris.New()
	app.SetName("[DEV-API]")
	app.Logger().SetOutput(os.Stdout)
	if api.DEBUG {
		app.Logger().SetLevel("debug")
	}else{
		app.Logger().SetLevel("info")
	}
	app.Logger().SetTimeFormat("2006-01-02 15:04:05.000:")
	app.Logger().Debugf("DEV API Debug mode on")
	app.Macros().Get("string").RegisterFunc("pid", api.PluginIdRe.MatchString)
	app.Macros().Get("string").RegisterFunc("version", api.VersionRe.MatchString)

	if !api.DEBUG {
		app.Use(recover.New())
	}
	app.Use(func(ctx iris.Context){
		var (
			ip, method, path string
			startTime time.Time
			usedTime time.Duration
			status int
		)
		startTime = time.Now()

		ip = ctx.RemoteAddr()
		if rip := ctx.GetHeader("X-Real-IP"); len(rip) > 0 {
			ip = rip
		}
		method = ctx.Method()
		path = ctx.Path()

		ctx.Next()

		usedTime = time.Since(startTime)
		status = ctx.GetStatusCode()
		line := fmt.Sprintf("%v %4v %s %s %s", status, usedTime, ip, method, path)

		if irisContext.StatusCodeNotSuccessful(ctx.GetStatusCode()) {
			ctx.Application().Logger().Warn(line)
		}else{
			ctx.Application().Logger().Debug(line)
		}
	})
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context){
		if !ctx.IsStopped() {
			ctx.StopWithJSON(iris.StatusNotFound, ErrResp{
				Status: "error",
				Name: "EntryPointNotExist",
				Msg: "The entry point you are accessing is not exists",
				Extra: ctx.Path(),
			})
			return
		}
	})

	app.Get("/", func(ctx iris.Context){
		ctx.JSON(iris.Map{
			"status": "ok",
			"time": time.Now().UTC().String(),
			"version": DevApiVerion,
		})
	})

	app.Get("/plugins", devPlugins)
	app.PartyFunc("/plugins", func(p iris.Party){
		p.Get("/", devPlugins)
		p.Get("/count", devPluginCounts)
	})
	app.PartyFunc("/plugin/{id:string pid()}", func(p iris.Party){
		p.Get("/info", devPluginInfo)
		p.Get("/readme", devPluginReadme)
		p.Get("/releases", devPluginReleases)
		p.PartyFunc("/release/{tag:string version()}", func(p iris.Party){
			p.Get("/", devPluginRelease)
			p.Get("/asset/{filename:file}", devPluginAsset)
		})
	})

	if err := app.Build(); err != nil {
		app.Logger().Fatalf("Cannot build application's router: %v", err)
	}

	server := &http.Server{
		Handler: app,
		ReadTimeout: time.Second * 30,
		WriteTimeout: time.Second * 60,
	}

	exit := make(chan struct{}, 0)

	go func(){
		defer close(exit)
		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			// kill -SIGINT XXXX or Ctrl+c
			os.Interrupt,
			syscall.SIGINT, // register that too, it should be ok
			// os.Kill  is equivalent with the syscall.Kill
			os.Kill,
			syscall.SIGKILL, // register that too, it should be ok
			// kill -SIGTERM XXXX
			syscall.SIGTERM,
		)
		select {
		case <-ch:
			timeout := 5 * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			server.Shutdown(ctx)
		}
	}()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		app.Logger().Fatalf("Error when server listening: %v", err)
	}

	app.Logger().Infof("Server listening at %s", listener.Addr().String())

	if err := server.Serve(listener); err != nil && err != iris.ErrServerClosed {
		app.Logger().Fatalf("Error when server running: %v", err)
	}
	select {
	case <-exit:
	case <-time.After(6 * time.Second):
		app.Logger().Warnf("Program exited incorrectly")
	}
}

func devPlugins(ctx iris.Context){
	var payload api.PluginListOpt
	if body, err := ctx.GetBody(); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, NewErrResp("BodyReadErr", err))
		return
	}else if len(body) > 0 {
		if err = json.Unmarshal(body, &payload); err != nil {
			ctx.StopWithJSON(iris.StatusBadRequest, NewErrResp("JsonDecodeErr", err))
			return
		}
	}
	if ctx.URLParamExists("filterBy") {
		payload.FilterBy = ctx.URLParamTrim("filterBy")
	}
	if tags0 := ctx.URLParamTrim("tags"); len(tags0) > 0 {
		payload.Tags = strings.Split(tags0, ",")
	}
	if ctx.URLParamExists("sortBy") {
		payload.SortBy = ctx.URLParamTrim("sortBy")
	}
	if ctx.URLParamExists("reversed") {
		payload.Reversed, _ = ctx.URLParamBool("reversed")
	}
	if ctx.URLParamExists("offset") {
		payload.Offset, _ = ctx.URLParamInt("offset")
	}
	if ctx.URLParamExists("limit") {
		payload.Limit, _ = ctx.URLParamInt("limit")
	}
	list, err := apiIns.GetPluginList(payload)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, NewErrResp("ApiErr", err))
		return
	}
	err = ctx.JSON(iris.Map{
		"status": "ok",
		"data": list,
	})
}

func devPluginCounts(ctx iris.Context){
	var payload api.PluginListOpt
	if body, err := ctx.GetBody(); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, NewErrResp("BodyReadErr", err))
		return
	}else if len(body) > 0 {
		if err = json.Unmarshal(body, &payload); err != nil {
			ctx.StopWithJSON(iris.StatusBadRequest, NewErrResp("JsonDecodeErr", err))
			return
		}
	}
	if ctx.URLParamExists("filterBy") {
		payload.FilterBy = ctx.URLParamTrim("filterBy")
	}
	if tags0 := ctx.URLParamTrim("tags"); len(tags0) > 0 {
		payload.Tags = strings.Split(tags0, ",")
	}
	if ctx.URLParamExists("sortBy") {
		payload.SortBy = ctx.URLParamTrim("sortBy")
	}
	if ctx.URLParamExists("reversed") {
		payload.Reversed, _ = ctx.URLParamBool("reversed")
	}
	if ctx.URLParamExists("offset") {
		payload.Offset, _ = ctx.URLParamInt("offset")
	}
	if ctx.URLParamExists("limit") {
		payload.Limit, _ = ctx.URLParamInt("limit")
	}
	counts, err := apiIns.GetPluginCounts(payload)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, NewErrResp("ApiErr", err))
		return
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"data": counts,
	})
}

func devPluginInfo(ctx iris.Context){
	id := ctx.Params().GetString("id")
	info, err := apiIns.GetPluginInfo(id)
	if err != nil {
		if err == api.ErrNotFound {
			ctx.StopWithJSON(iris.StatusNotFound, NewErrResp("NotFound", err))
			return
		}
		ctx.StopWithJSON(iris.StatusInternalServerError, NewErrResp("ApiErr", err))
		return
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"data": info,
	})
}

func devPluginReadme(ctx iris.Context){
	id := ctx.Params().GetString("id")
	render, _ := ctx.URLParamBool("render")
	content, err := apiIns.GetPluginReadme(id)
	if err != nil {
		if err == api.ErrNotFound {
			ctx.StopWithJSON(iris.StatusNotFound, NewErrResp("NotFound", err))
			return
		}
		ctx.StopWithJSON(iris.StatusInternalServerError, NewErrResp("ApiErr", err))
		return
	}
	body := content.Data
	if render {
		body = api.RenderMarkdown(body, &api.Option{
			URLPrefix: content.URLPrefix,
			DataURLPrefix: content.DataURLPrefix,
			HeadingIDPrefix: "~mdh@",
		})
	}
	_, _ = ctx.Write(body)
}

func devPluginReleases(ctx iris.Context){
	id := ctx.Params().GetString("id")
	releases, err := apiIns.GetPluginReleases(id)
	if err != nil {
		if err == api.ErrNotFound {
			ctx.StopWithJSON(iris.StatusNotFound, NewErrResp("NotFound", err))
			return
		}
		ctx.StopWithJSON(iris.StatusInternalServerError, NewErrResp("ApiErr", err))
		return
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"data": releases,
	})
}

func devPluginRelease(ctx iris.Context){
	id := ctx.Params().GetString("id")
	tag, err := api.VersionFromString(ctx.Params().GetString("tag"))
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, NewErrResp("VersionFormatErr", err))
		return
	}
	release, err := apiIns.GetPluginRelease(id, tag)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, NewErrResp("ApiErr", err))
		return
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"data": release,
	})
}

func devPluginAsset(ctx iris.Context){
	id := ctx.Params().GetString("id")
	tag, err := api.VersionFromString(ctx.Params().GetString("tag"))
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, NewErrResp("VersionFormatErr", err))
		return
	}
	filename := ctx.Params().GetString("filename")
	fd, modTime, err := apiIns.GetPluginReleaseAsset(id, tag, filename)
	if err != nil {
		if err == api.ErrNotFound {
			ctx.StopWithJSON(iris.StatusBadRequest, NewErrResp("NotFound", err))
			return
		}
		ctx.StopWithJSON(iris.StatusInternalServerError, NewErrResp("ApiErr", err))
		return
	}
	defer fd.Close()
	ctx.ServeContent(fd, filename, modTime)
}
