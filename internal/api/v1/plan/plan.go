package plan

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/hantbk/vbackup/internal/cron"
	"github.com/hantbk/vbackup/internal/entity/v1/plan"
	"github.com/hantbk/vbackup/internal/model"
	"github.com/hantbk/vbackup/internal/server"
	"github.com/hantbk/vbackup/internal/service/v1/common"
	ser "github.com/hantbk/vbackup/internal/service/v1/plan"
	"github.com/hantbk/vbackup/pkg/utils"
)

var planServer ser.Service

func init() {
	planServer = ser.GetService()
}

func createHandler() iris.Handler {
	return func(ctx *context.Context) {
		var p plan.Plan
		err := ctx.ReadJSON(&p)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		if p.Path == "" {
			utils.ErrorStr(ctx, "Path cannot be empty")
			return
		}
		if p.Status == -1 {
			p.Status = plan.StopStatus
		}
		if p.ExecTimeCron == "" {
			utils.ErrorStr(ctx, "Cron expression cannot be empty")
			return
		}
		_, err = cron.GetNextTimes(p.ExecTimeCron)
		if err != nil {
			utils.ErrorStr(ctx, "Invalid cron expression format: "+err.Error())
			return
		}
		err = planServer.Create(&p, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		if p.Status == plan.RunningStatus {
			err := cron.AddJob(p.ExecTimeCron, cron.BackupJob{
				PlanId: p.Id,
			})
			if err != nil {
				utils.Errore(ctx, err)
				return
			}
		}
		ctx.Values().Set("data", p.Id)
	}
}

func deleteHandler() iris.Handler {
	return func(ctx *context.Context) {
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		err = planServer.Delete(id, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		cron.ClearJob()
		initPlan()
		ctx.Values().Set("data", "")
	}
}

func updateHandler() iris.Handler {
	return func(ctx *context.Context) {
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		var p plan.Plan
		err = ctx.ReadJSON(&p)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		if p.Path == "" {
			utils.ErrorStr(ctx, "path can not be empty")
			return
		}
		if p.Status == -1 {
			p.Status = plan.StopStatus
		}
		_, err = cron.GetNextTimes(p.ExecTimeCron)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		p.Id = id
		err = planServer.Update(&p, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		cron.ClearJob()
		initPlan()
		ctx.Values().Set("data", "")
	}
}

func searchHandler() iris.Handler {
	return func(ctx *context.Context) {
		res := model.PageParam(ctx)
		status, err := ctx.URLParamInt("status")
		if err != nil {
			status = 0
		}
		repositoryId, err := ctx.URLParamInt("repositoryId")
		if err != nil {
			repositoryId = 0
		}

		path := ctx.URLParam("path")
		name := ctx.URLParam("name")
		total, plans, err := planServer.Search(res.PageNum, res.PageSize, status, repositoryId, path, name, common.DBOptions{})
		if err != nil && err.Error() != "not found" {
			utils.Errore(ctx, err)
			return
		}
		res.Total = total
		res.Items = plans
		ctx.Values().Set("data", res)
	}
}
func getNextTime() iris.Handler {
	return func(ctx *context.Context) {
		cronStr := ctx.URLParam("cron")
		res, err := cron.GetNextTimes(cronStr)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", res)
	}
}

func Install(parent iris.Party) {
	// Plan-related APIs
	planParty := parent.Party("/plan")
	// Create
	planParty.Post("", createHandler())
	// Delete
	planParty.Delete("/:id", deleteHandler())
	// Update
	planParty.Put("/:id", updateHandler())
	// Search
	planParty.Get("", searchHandler())
	planParty.Get("/next_time", getNextTime())
	initPlan()
}

// initPlan Initializes plans into scheduled tasks
func initPlan() {
	plans, err := planServer.List(plan.RunningStatus, common.DBOptions{})
	if err != nil {
		fmt.Println("No scheduled tasks loaded")
		return
	}
	for _, p := range plans {
		if p.ExecTimeCron == "" {
			continue
		}
		err = cron.AddJob(p.ExecTimeCron, cron.BackupJob{
			PlanId: p.Id,
		})
		if err != nil {
			server.Logger().Error(err)
		}
	}
	fmt.Println("Scheduled tasks loaded successfully")
}
