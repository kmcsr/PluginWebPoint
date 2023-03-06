
package main

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
)

var (
	PluginIdRe = regexp.MustCompile("[0-9a-z_]{1,64}")
)

func devPluginList(ctx iris.Context){
	loger.Debugf("URL params: %v", ctx.URLParams())
	filterBy := ctx.URLParamTrim("filterBy")
	tags0 := ctx.URLParamTrim("tags")
	var tags []string
	if len(tags0) > 0 {
		tags = strings.Split(tags0, ",")
	}
	sortBy := ctx.URLParamTrim("sortBy")
	reversed, _ := ctx.URLParamBool("reversed")
	list, err := APIIns.GetPluginList(filterBy, tags, sortBy, reversed)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"status": "error",
			"error": "apiError",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"data": list,
	})
}

func devPluginInfo(ctx iris.Context){
	id := ctx.Params().GetString("id")
	info, err := APIIns.GetPluginInfo(id)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"status": "error",
			"error": "apiError",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"data": info,
	})
}

func GetDevAPIHandler()(http.Handler){
	app := iris.New()
	app.SetName("dev-api")
	app.Macros().Get("string").RegisterFunc("pid", PluginIdRe.MatchString)

	app.Get("/", func(ctx iris.Context){
		ctx.JSON(iris.Map{
			"status": "ok",
			"time": time.Now().UTC().String(),
		})
	})

	app.PartyFunc("/plugin", func(p iris.Party){
		p.Get("/list", devPluginList)
		p.PartyFunc("/{id:string pid()}", func(p iris.Party){
			p.Get("/info", devPluginInfo)
		})
	})


	if err := app.Build(); err != nil {
		loger.Panicf("Cannot build router: %v", err)
		return nil
	}
	return app
}
