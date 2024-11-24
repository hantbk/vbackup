package repository

import (
	"strconv"

	"github.com/hantbk/vbackup/internal/entity/v1/repository"
	"github.com/hantbk/vbackup/internal/service/v1/common"
	policyDao "github.com/hantbk/vbackup/internal/service/v1/policy"
	repositoryDao "github.com/hantbk/vbackup/internal/service/v1/repository"
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
		var rep repository.Repository
		err := ctx.ReadJSON(&rep)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		if rep.Password == "" {
			utils.ErrorStr(ctx, "Please enter the password")
			return
		}
		option, _ := resticProxy.GetGlobalOptions(rep)
		repo, err1 := resticProxy.OpenRepository(ctx, option)
		if err1 != nil {
			// Repository error, reinitialize
			version, err := resticProxy.RunInit(ctx, option)
			if err != nil {
				utils.Errore(ctx, err)
				return
			}
			rep.RepositoryVersion = strconv.Itoa(int(version))
		} else {
			rep.RepositoryVersion = strconv.Itoa(int(repo.Config().Version))
		}
		rep.Status = repository.StatusNone
		err = repositoryService.Create(&rep, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		go resticProxy.InitRepository()
		ctx.Values().Set("data", rep.Id)
	}
}

func listHandler() iris.Handler {
	return func(ctx *context.Context) {
		repotype, err := ctx.URLParamInt("type")
		if err != nil {
			repotype = 0
		}
		name := ctx.URLParam("name")
		ress, err := resticProxy.GetAllRepoWithStatus(repotype, name)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", ress)
	}
}

func delHanlder() iris.Handler {
	return func(ctx *context.Context) {
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		err = repositoryService.Delete(id, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		_ = policyService.DeleteByRepo(id, common.DBOptions{})
		ctx.Values().Set("data", "")
	}
}

func updateHandler() iris.Handler {
	return func(ctx *context.Context) {
		var rep repository.Repository
		err := ctx.ReadJSON(&rep)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		// Only the following fields are allowed to be modified
		rep2, err := repositoryService.Get(id, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		if rep.Name != "" {
			rep2.Name = rep.Name
		}
		if rep.KeyId != "" {
			rep2.KeyId = rep.KeyId
		}
		if rep.Region != "" {
			rep2.Region = rep.Region
		}
		if rep.Bucket != "" {
			rep2.Bucket = rep.Bucket
		}
		if rep.Secret != "" {
			rep2.Secret = rep.Secret
		}
		if rep.Endpoint != "" {
			rep2.Endpoint = rep.Endpoint
		}
		rep2.PackSize = rep.PackSize
		option, _ := resticProxy.GetGlobalOptions(*rep2)
		_, err = resticProxy.OpenRepository(ctx, option)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		err = repositoryService.Update(rep2, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		go resticProxy.InitRepository()
		ctx.Values().Set("data", "")
	}
}

func getHandler() iris.Handler {
	return func(ctx *context.Context) {

		id, err := ctx.Params().GetInt("id")
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		resp, err := repositoryService.Get(id, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}

		config := resticProxy.CheckRepoStatus(resp.Id)
		if config != nil {
			resp.Status = repository.StatusRun
			resp.RepositoryVersion = strconv.Itoa(int(config.Version))
		} else {
			resp.Status = repository.StatusErr
			resp.Errmsg = "Repository connection timeout"

		}
		resp.Password = "******"
		ctx.Values().Set("data", resp)
	}
}

// func errHandler() iris.Handler {
// 	return func(ctx *context.Context) {
// 		err := ctx.GetErr()
// 		if err != nil {
// 			utils.Errore(ctx, err)
// 		}
// 	}
// }

func Install(parent iris.Party) {
	// Repository related endpoints
	sp := parent.Party("/repository")
	// Create new repository
	sp.Post("", createHandler())
	// List all repositories
	sp.Get("", listHandler())
	// Delete repository by ID
	sp.Delete("/:id", delHanlder())
	// Update repository by ID
	sp.Put("/:id", updateHandler())
	// Get repository by ID
	sp.Get("/:id", getHandler())
	// Error handler
	// sp.Use(errHandler())
}
