
package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	irisLogger "github.com/kataras/iris/v12/middleware/logger"
)

const DevApiVerion = 0

type ErrResp struct{
	Status string `json:"status"`
	Name   string `json:"error"`
	Msg    string `json:"message"`
}

func NewErrResp(name string, err error)(*ErrResp){
	return &ErrResp{
		Status: "error",
		Name: name,
		Msg: err.Error(),
	}
}

func devPlugins(ctx iris.Context){
	var payload PluginListOpt
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
	list, err := Ins.GetPluginList(payload)
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
	var payload PluginListOpt
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
	counts, err := Ins.GetPluginCounts(payload)
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
	info, err := Ins.GetPluginInfo(id)
	if err != nil {
		if err == ErrNotFound {
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
	body, prefix, err := Ins.GetPluginReadme(id)
	if err != nil {
		if err == ErrNotFound {
			ctx.StopWithJSON(iris.StatusNotFound, NewErrResp("NotFound", err))
			return
		}
		ctx.StopWithJSON(iris.StatusInternalServerError, NewErrResp("ApiErr", err))
		return
	}
	if render {
		body = RenderMarkdown(body, &Option{
			URLPrefix: prefix,
		})
	}
	_, _ = ctx.Write(body)
}

func devPluginReleases(ctx iris.Context){
	id := ctx.Params().GetString("id")
	releases, err := Ins.GetPluginReleases(id)
	if err != nil {
		if err == ErrNotFound {
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
	tag, err := VersionFromString(ctx.Params().GetString("tag"))
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, NewErrResp("VersionFormatErr", err))
		return
	}
	release, err := Ins.GetPluginRelease(id, tag)
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
	tag, err := VersionFromString(ctx.Params().GetString("tag"))
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, NewErrResp("VersionFormatErr", err))
		return
	}
	filename := ctx.Params().GetString("filename")
	fd, modTime, err := Ins.GetPluginReleaseAsset(id, tag, filename)
	if err != nil {
		if err == ErrNotFound {
			ctx.StopWithJSON(iris.StatusBadRequest, NewErrResp("NotFound", err))
			return
		}
		ctx.StopWithJSON(iris.StatusInternalServerError, NewErrResp("ApiErr", err))
		return
	}
	defer fd.Close()
	ctx.ServeContent(fd, filename, modTime)
}

func GetDevAPIHandler()(http.Handler){
	app := iris.New()
	app.SetName("dev-api")
	app.Macros().Get("string").RegisterFunc("pid", PluginIdRe.MatchString)
	app.Macros().Get("string").RegisterFunc("version", VersionRe.MatchString)

	app.Use(recover.New())
	app.Use(irisLogger.New())
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context){
		if !ctx.IsStopped() {
			ctx.StopWithJSON(iris.StatusNotFound, ErrResp{
				Status: "error",
				Name: "EntryPointNotExist",
				Msg: "The entry point you are accessing is not exists",
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
		loger.Panicf("Cannot build router: %v", err)
		return nil
	}
	return app
}
