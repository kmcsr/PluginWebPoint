
package api

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	irisLogger "github.com/kataras/iris/v12/middleware/logger"
)

func devPlugins(ctx iris.Context){
	loger.Debugf("URL params: %v", ctx.URLParams())
	filterBy := ctx.URLParamTrim("filterBy")
	tags0 := ctx.URLParamTrim("tags")
	var tags []string
	if len(tags0) > 0 {
		tags = strings.Split(tags0, ",")
	}
	sortBy := ctx.URLParamTrim("sortBy")
	reversed, _ := ctx.URLParamBool("reversed")
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	list, err := Ins.GetPluginList(PluginListOpt{
		FilterBy: filterBy,
		Tags: tags,
		SortBy: sortBy,
		Reversed: reversed,
		Offset: offset,
		Limit: limit,
	})
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"status": "error",
			"error": "ApiError",
			"message": err.Error(),
		})
		return
	}
	err = ctx.JSON(iris.Map{
		"status": "ok",
		"data": list,
	})
}

func devPluginCounts(ctx iris.Context){
	filterBy := ctx.URLParamTrim("filterBy")
	sortBy := ctx.URLParamTrim("sortBy")
	reversed, _ := ctx.URLParamBool("reversed")
	offset, _ := ctx.URLParamInt("offset")
	limit, _ := ctx.URLParamInt("limit")
	counts, err := Ins.GetPluginCounts(PluginListOpt{
		FilterBy: filterBy,
		SortBy: sortBy,
		Reversed: reversed,
		Offset: offset,
		Limit: limit,
	})
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"status": "error",
			"error": "ApiError",
			"message": err.Error(),
		})
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
			ctx.StopWithJSON(iris.StatusNotFound, iris.Map{
				"status": "error",
				"error": "NotFound",
				"message": err.Error(),
			})
			return
		}
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"status": "error",
			"error": "ApiError",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(iris.Map{
		"status": "ok",
		"data": info,
	})
}

func devPluginReleases(ctx iris.Context){
	id := ctx.Params().GetString("id")
	releases, err := Ins.GetPluginReleases(id)
	if err != nil {
		if err == ErrNotFound {
			ctx.StopWithJSON(iris.StatusNotFound, iris.Map{
				"status": "error",
				"error": "NotFound",
				"message": err.Error(),
			})
			return
		}
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"status": "error",
			"error": "ApiError",
			"message": err.Error(),
		})
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
		if err == ErrNotFound {
			ctx.StopWithJSON(iris.StatusNotFound, iris.Map{
				"status": "error",
				"error": "NotFound",
				"message": err.Error(),
			})
			return
		}
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"status": "error",
			"error": "VersionFormatError",
			"message": err.Error(),
		})
		return
	}
	release, err := Ins.GetPluginRelease(id, tag)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"status": "error",
			"error": "ApiError",
			"message": err.Error(),
		})
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
		if err == ErrNotFound {
			ctx.StopWithJSON(iris.StatusNotFound, iris.Map{
				"status": "error",
				"error": "NotFound",
				"message": err.Error(),
			})
			return
		}
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"status": "error",
			"error": "VersionFormatError",
			"message": err.Error(),
		})
		return
	}
	filename := ctx.Params().GetString("filename")
	fd, modTime, err := Ins.GetPluginReleaseAsset(id, tag, filename)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.StopWithJSON(iris.StatusNotFound, iris.Map{
				"status": "error",
				"error": "FileNotFoundError",
				"message": err.Error(),
			})
			return
		}
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"status": "error",
			"error": "ApiError",
			"message": err.Error(),
		})
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

	app.Get("/", func(ctx iris.Context){
		ctx.JSON(iris.Map{
			"status": "ok",
			"time": time.Now().UTC().String(),
		})
	})

	app.Get("/plugins", devPlugins)
	app.PartyFunc("/plugins", func(p iris.Party){
		p.Get("/", devPlugins)
		p.Get("/count", devPluginCounts)
	})
	app.PartyFunc("/plugin", func(p iris.Party){
		p.PartyFunc("/{id:string pid()}", func(p iris.Party){
			p.Get("/info", devPluginInfo)
			p.Get("/releases", devPluginReleases)
			p.PartyFunc("/release/{tag:string version()}", func(p iris.Party){
				p.Get("/", devPluginRelease)
				p.Get("/asset/{filename:file}", devPluginAsset)
			})
		})
	})


	if err := app.Build(); err != nil {
		loger.Panicf("Cannot build router: %v", err)
		return nil
	}
	return app
}
