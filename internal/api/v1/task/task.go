package task

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hantbk/vbackup/internal/consts"
	thmodel "github.com/hantbk/vbackup/internal/entity/v1/task"
	"github.com/hantbk/vbackup/internal/model"
	"github.com/hantbk/vbackup/internal/server"
	"github.com/hantbk/vbackup/internal/service/v1/common"
	"github.com/hantbk/vbackup/internal/service/v1/plan"
	ser "github.com/hantbk/vbackup/internal/service/v1/task"
	"github.com/hantbk/vbackup/internal/store/task"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/restic"
	"github.com/hantbk/vbackup/pkg/utils"
	resticProxy "github.com/hantbk/vbackup/restic_proxy"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

var taskService ser.Service
var planService plan.Service

func init() {
	taskService = ser.GetService()
	planService = plan.GetService()
}

func backupHandler() iris.Handler {
	return func(ctx *context.Context) {
		planid, err := ctx.Params().GetInt("planid")
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		var taskid int
		taskid, err = Backup(planid)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", taskid)
	}
}

// Backup - Backup data, planid is the plan ID
func Backup(planid int) (int, error) {
	pl, err := planService.Get(planid, common.DBOptions{})
	if err != nil {
		return 0, err
	}
	repoid := pl.RepositoryId
	progress := &model.StatusUpdate{
		MessageType:      "status",
		SecondsElapsed:   "0",
		SecondsRemaining: "0",
		TotalFiles:       0,
		FilesDone:        0,
		TotalBytes:       "0",
		BytesDone:        "0",
		ErrorCount:       0,
	}
	ta := &thmodel.Task{
		Path:            pl.Path,
		Name:            "backup_" + strconv.Itoa(planid) + "_" + strconv.Itoa(repoid) + "_" + time.Now().Format(consts.TaskHistoryName),
		Status:          task.StatusNew,
		RepositoryId:    repoid,
		Progress:        progress,
		PlanId:          pl.Id,
		ReadConcurrency: pl.ReadConcurrency,
	}
	err = taskService.Create(ta, common.DBOptions{})
	if err != nil {
		return 0, err
	}
	opt := resticProxy.BackupOptions{UseFsSnapshot: true, ReadConcurrency: pl.ReadConcurrency}
	taskInfo := task.TaskInfo{
		Name: ta.Name,
		Path: ta.Path,
	}
	taskInfo.SetId(ta.Id)

	err = resticProxy.RunBackup(opt, repoid, taskInfo)

	if err != nil {

		ta.ArchivalError = append(ta.ArchivalError, model.ErrorUpdate{
			MessageType: "error",
			Error:       err.Error(),
		})
		ta.Status = task.StatusError

		_ = taskService.Update(ta, common.DBOptions{})
		return 0, err
	}

	return taskInfo.GetId(), nil
}

func restoreHandler() iris.Handler {
	return func(ctx *context.Context) {
		snapshotid := ctx.Params().Get("snapshotid")
		repository, err := ctx.Params().GetInt("repository")
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		var info model.RestoreInfo
		err = ctx.ReadJSON(&info)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		target := info.Target
		if target == "" {
			target = string(filepath.Separator)
		}
		path := info.Paths
		host := info.Hosts
		tag := info.Tags
		exclude := info.Exclude
		iexclude := info.IExclude
		include := info.Include
		iinclude := info.IInclude

		var paths []string
		if path != "" {
			paths = strings.Split(path, ",")
		}
		var hosts []string
		if host != "" {
			hosts = strings.Split(host, ",")
		}
		var excludes []string
		if exclude != "" {
			excludes = strings.Split(exclude, ",")
		}
		var includes []string
		if include != "" {
			includes = strings.Split(include, ",")
		}
		var iexcludes []string
		if iexclude != "" {
			iexcludes = strings.Split(iexclude, ",")
		}
		var iincludes []string
		if iinclude != "" {
			iincludes = strings.Split(iinclude, ",")
		}
		tags := restic.TagLists{}
		if tag != "" {
			err := tags.Set(tag)
			if err != nil {
				utils.Errore(ctx, err)
				return
			}
		}
		opts := resticProxy.RestoreOptions{
			Exclude:            excludes,
			InsensitiveExclude: iexcludes,
			Include:            includes,
			InsensitiveInclude: iincludes,
			Target:             target,
			SnapshotFilter:     restic.SnapshotFilter{Hosts: hosts, Paths: paths, Tags: tags},
			Verify:             info.Verify,
		}

		err = resticProxy.RunRestore(opts, repository, snapshotid)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", "")
	}
}

func searchHandler() iris.Handler {
	return func(ctx *context.Context) {
		res := model.PageParam(ctx)
		status, err := ctx.URLParamInt("status")
		if err != nil {
			status = -1
		}
		repositoryId, err := ctx.URLParamInt("repositoryId")
		if err != nil {
			repositoryId = 0
		}
		planId, err := ctx.URLParamInt("planId")
		if err != nil {
			planId = 0
		}
		path := ctx.URLParam("path")
		name := ctx.URLParam("name")

		total, taskHistories, err := taskService.Search(
			res.PageNum, res.PageSize, status, repositoryId, planId, path, name, common.DBOptions{})
		if err != nil && err.Error() != "not found" {
			utils.Errore(ctx, err)
			return
		}
		res.Total = total
		for key, t := range taskHistories {
			if t.Status == task.StatusRunning {
				ta := task.TaskInfos.Get(t.Id)
				if ta != nil {
					t.Progress = ta.(*task.TaskInfo).Progress
				}
			}
			taskHistories[key] = t
		}
		res.Items = taskHistories
		ctx.Values().Set("data", res)
	}
}

// ClearTaskRunning - Clean abnormal tasks
func ClearTaskRunning() {
	server.Logger().Debugln("Start executing ClearTaskRunning")
	opt := common.DBOptions{}
	_, taskHistories, _ := taskService.Search(
		1, 10, -1, 0, 0, "", "", opt)
	errTasks := make([]thmodel.Task, 0)
	for _, t := range taskHistories {
		if t.Status == task.StatusRunning {
			ta := task.TaskInfos.Get(t.Id)
			if ta == nil && t.Summary == nil {
				t.Status = task.StatusError
				t.ArchivalError = append(t.ArchivalError, model.ErrorUpdate{
					MessageType: "error",
					Error:       "Unknown status exception",
				})
			}
			if t.Summary != nil {
				t.Status = task.StatusEnd
			}
			errTasks = append(errTasks, t)
		} else if t.Status == task.StatusError {
			if t.Summary != nil {
				t.Status = task.StatusEnd
			}
			errTasks = append(errTasks, t)
		}

	}
	for _, v := range errTasks {
		_ = taskService.Update(&v, opt)
	}
}

func Install(parent iris.Party) {
	// Task-related APIs
	taskParty := parent.Party("/task")
	// Add backup task
	taskParty.Post("/backup/:planid", backupHandler())
	// Add restore task
	taskParty.Post("/:repository/restore/:snapshotid/", restoreHandler())
	// Search tasks
	taskParty.Get("", searchHandler())
}
