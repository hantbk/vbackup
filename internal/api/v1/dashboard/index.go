package dashboard

import (
	"github.com/hantbk/vbackup/internal/consts"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

var planServer ser.Service
var logServer logser.Service

func init() {
	planServer = ser.GetService()
	logServer = logser.GetService()
}

func indexHandler() iris.Handler {
	return func(ctx *context.Context) {
		var (
			plant, planc, repot, repoc int
		)
		plans, err := planServer.List(-1, common.DBOptions{})
		if err != nil {
			plant = 0
			planc = 0
		}
		plant, planc = getPlanData(plans)

		planinfo := model.PlanInfo{
			Total:        plant,
			RunningCount: planc,
		}
		repositories, err := resticProxy.GetAllRepoWithStatus(0, "")
		if err != nil {
			repot = 0
			repoc = 0
		}
		key1 := consts.Key("GetAllRepoStats", "backupinfo")
		key2 := consts.Key("GetAllRepoStats", "backupinfos")
		c := server.Cache()
		repot, repoc = getRepoRunCount(repositories)

		repositorieinfo := model.RepositoryInfo{
			Total:        repot,
			RunningCount: repoc,
		}
		backupinfo := model.BackupInfo{
			FileTotal:    0,
			DataDay:      "0",
			DataSize:     0,
			DataSizeStr:  "0",
			SnapshotsNum: 0,
		}
		backupinfos := make([]model.BackupInfo, 0)
		biold, has := c.Get(key1)
		if has {
			backupinfo = biold.(model.BackupInfo)
			biolds, has2 := c.Get(key2)
			if has2 {
				backupinfos = biolds.([]model.BackupInfo)
			} else {
				go resticProxy.GetAllRepoStats()
			}
		} else {
			go resticProxy.GetAllRepoStats()
		}

		boardinfo := model.BoardInfo{
			PlanInfo:       planinfo,
			RepositoryInfo: repositorieinfo,
			BackupInfo:     backupinfo,
			BackupInfos:    backupinfos,
		}
		ctx.Values().Set("data", boardinfo)
	}
}

func searchLogHandler() iris.Handler {
	return func(ctx *context.Context) {
		res := model.PageParam(ctx)
		operator := ctx.URLParam("operator")
		operation := ctx.URLParam("operation")
		url := ctx.URLParam("url")
		data := ctx.URLParam("data")
		total, logs, err := logServer.Search(res.PageNum, res.PageSize, operator, operation, url, data, common.DBOptions{})
		if err != nil && err.Error() != "not found" {
			utils.Errore(ctx, err)
			return
		}
		res.Total = total
		res.Items = logs
		ctx.Values().Set("data", res)
	}
}

func doGetAllRepoStatsHandler() iris.Handler {
	return func(ctx *context.Context) {
		go resticProxy.GetAllRepoStats()
		ctx.Values().Set("data", "")
	}
}
func Install(parent iris.Party) {
	dashboardParty := parent.Party("/dashboard")
	// Home Page Statistics Data
	dashboardParty.Get("/index", indexHandler())
	dashboardParty.Post("/doGetAllRepoStats", doGetAllRepoStatsHandler())
	dashboardParty.Get("/logs", searchLogHandler())
}

func getPlanData(plans []plan.Plan) (total, runc int) {
	total = len(plans)
	runc = 0
	for _, p := range plans {
		if p.Status == plan.RunningStatus {
			runc++
		}
	}
	return
}

func getRepoRunCount(repositories []repository.Repository) (total, runc int) {
	total = len(repositories)
	runc = 0
	for _, v := range repositories {
		if v.Status == repository.StatusRun {
			runc++
		}
	}
	return
}
