package system

import (
	"encoding/json"

	"github.com/hantbk/vbackup"
	"github.com/hantbk/vbackup/internal/consts/global"
	sys "github.com/hantbk/vbackup/internal/entity/v1/system"
	"github.com/hantbk/vbackup/internal/service/v1/system"
	fileutil "github.com/hantbk/vbackup/pkg/file"
	"github.com/hantbk/vbackup/pkg/utils"
	"github.com/hantbk/vbackup/pkg/utils/http"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func lsHandler() iris.Handler {
	return func(ctx *context.Context) {
		path := ctx.URLParam("path")
		listDir, err := fileutil.ListDir(path)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", listDir)
	}
}

func upgradeVersionHandler() iris.Handler {
	return func(ctx *context.Context) {
		version := ctx.Params().GetString("version")
		v := vbackup.GetVersion()
		if v.Version == version {
			utils.ErrorStr(ctx, "The update version is the same as the current version")
		}
		err := system.Upgrade(version)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", "upgrading")
	}
}

func versionHandler() iris.Handler {
	return func(ctx *context.Context) {
		v := vbackup.GetVersion()
		ctx.Values().Set("data", v)
	}
}

func latestVersionHandler() iris.Handler {
	return func(ctx *context.Context) {
		body, err := http.Get(global.LatestUrl)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		var releases []sys.Release
		err = json.Unmarshal([]byte(body), &releases)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		if len(releases) == 0 {
			utils.ErrorStr(ctx, "Failed to retrieve the new version")
		}
		latest := releases[0].TagName
		ctx.Values().Set("data", latest)
	}
}

func userHomePathHandler() iris.Handler {
	return func(ctx *context.Context) {
		// Gọi hàm GetUserHomePath từ utils
		homePath := utils.GetUserHomePath() 
		if homePath == "" {
			utils.ErrorStr(ctx, "Failed to retrieve the user's home directory")
			return
		}
		ctx.Values().Set("data", homePath)
	}
}

func Install(parent iris.Party) {
	// System interfaces
	sp := parent.Party("/system")
	// List folders
	sp.Get("/ls", lsHandler())

	sp.Post("/upgradeVersion/:version", upgradeVersionHandler())

	sp.Get("/version", versionHandler())

	sp.Get("/version/latest", latestVersionHandler())

	// User home path
	sp.Get("/userHome", userHomePathHandler())
}
