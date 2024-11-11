package cron

import (
	"github.com/hantbk/vbackup/internal/api/v1/task"
	"github.com/hantbk/vbackup/internal/server"
)

type BackupJob struct {
	PlanId int
}

func (b BackupJob) Run() {
	_, err := task.Backup(b.PlanId)
	if err != nil {
		server.Logger().Error(err)
		return
	}
}

func (b BackupJob) GetType() int {
	return JOB_TYPE_BACKUP
}

var _ BaseJob = &BackupJob{}
