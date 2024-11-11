package policy

import (
	"github.com/hantbk/vbackup/internal/entity/v1/repository"
	"github.com/hantbk/vbackup/internal/server"
	"github.com/hantbk/vbackup/internal/service/v1/common"
	policyDao "github.com/hantbk/vbackup/internal/service/v1/policy"
	repositoryDao "github.com/hantbk/vbackup/internal/service/v1/repository"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/restic"
	"github.com/hantbk/vbackup/pkg/utils"
	resticProxy "github.com/hantbk/vbackup/restic_proxy"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

var policyService policyDao.Service
var repositoryService repositoryDao.Service

func init() {
	policyService = policyDao.GetService()
	repositoryService = repositoryDao.GetService()
}

func createHandler() iris.Handler {
	return func(ctx *context.Context) {
		var policy repository.ForgetPolicy
		err := ctx.ReadJSON(&policy)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		if policy.RepositoryId <= 0 {
			utils.ErrorStr(ctx, "Repository ID cannot be empty")
			return
		}
		if policy.Path == "" {
			utils.ErrorStr(ctx, "Path cannot be empty")
			return
		}
		err = policyService.Create(&policy, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", policy.Id)
	}
}

func updateHandler() iris.Handler {
	return func(ctx *context.Context) {
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		var policy repository.ForgetPolicy
		err = ctx.ReadJSON(&policy)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		if policy.RepositoryId <= 0 {
			utils.ErrorStr(ctx, "Repository ID cannot be empty")
			return
		}
		forgetPolicy, err := policyService.Get(id, common.DBOptions{})
		if err != nil {
			return
		}
		forgetPolicy.Value = policy.Value
		forgetPolicy.Type = policy.Type
		forgetPolicy.RepositoryId = policy.RepositoryId
		err = policyService.Update(forgetPolicy, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", policy.Id)
	}
}

func listHandler() iris.Handler {
	return func(ctx *context.Context) {
		repoid, err := ctx.URLParamInt("repository")
		if err != nil {
			repoid = 0
		}
		path := ctx.URLParam("path")
		policies, err := policyService.Search(repoid, path, common.DBOptions{})
		if err != nil && err.Error() != "not found" {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", policies)
	}
}

func delHanlder() iris.Handler {
	return func(ctx *context.Context) {
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		err = policyService.Delete(id, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", "")
	}
}

func doHanlder() iris.Handler {
	return func(ctx *context.Context) {
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		policy, err := policyService.Get(id, common.DBOptions{})
		if err != nil {
			return
		}
		opt := resticProxy.ForgetOptions{
			Prune:          true,
			SnapshotFilter: restic.SnapshotFilter{Paths: []string{policy.Path}},
		}
		setType(policy.Type, policy.Value, &opt)
		operid, err := resticProxy.RunForget(opt, policy.RepositoryId, []string{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", operid)
	}
}

func Install(parent iris.Party) {
	// Repository-related APIs
	sp := parent.Party("/policy")
	// Add new
	sp.Post("", createHandler())
	// Update
	sp.Put("/:id", updateHandler())
	// List
	sp.Get("", listHandler())
	// Delete
	sp.Delete("/:id", delHanlder())
	// Execute immediately
	sp.Post("/do/:id", doHanlder())
}

func DoPolicy() {
	reps, err := repositoryService.List(0, "", common.DBOptions{})
	if err != nil {
		return
	}
	for _, rep := range reps {
		policys, err := policyService.Search(rep.Id, "", common.DBOptions{})
		if err != nil {
			return
		}
		for i, policy := range policys {
			opt := resticProxy.ForgetOptions{
				Prune:          i == (len(policys) - 1),
				SnapshotFilter: restic.SnapshotFilter{Paths: []string{policy.Path}},
			}
			setType(policy.Type, policy.Value, &opt)
			err = resticProxy.RunForgetSync(opt, policy.RepositoryId, []string{})
			if err != nil {
				server.Logger().Error(err)
			}
			server.Logger().Infof("Cleaning %s under %s, keeping the latest %d %s snapshots", rep.Name, policy.Path, policy.Value, policy.Type)
		}
	}
}

func setType(t string, value int, opt *resticProxy.ForgetOptions) {
	v := resticProxy.ForgetPolicyCount(value)
	switch t {
	case "last":
		opt.Last = v
		break
	case "hourly":
		opt.Hourly = v
		break
	case "daily":
		opt.Daily = v
		break
	case "weekly":
		opt.Weekly = v
		break
	case "monthly":
		opt.Monthly = v
		break
	case "yearly":
		opt.Yearly = v
		break
	}
}
