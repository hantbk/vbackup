package cron

import "github.com/robfig/cron/v3"

type BaseJob interface {
	GetType() int
	cron.Job
}

const (
	JOB_TYPE_SYSTEM = 0 // System task, will not be deleted
	JOB_TYPE_BACKUP = 1 // Backup task
	JOB_TYPE_PRUNE  = 2 // Cleanup task
)
