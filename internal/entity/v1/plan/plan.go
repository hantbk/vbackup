package plan

import "github.com/hantbk/vbackup/internal/entity/v1/common"

type Plan struct {
	common.BaseModel `storm:"inline"`
	Name             string `json:"name"`
	Path             string `json:"path"`            // Backup or restore path
	RepositoryId     int    `json:"repositoryId"`    // ID of the repository
	Status           int    `json:"status"`          // Current status of the plan
	ExecTimeCron     string `json:"execTimeCron"`    // Scheduled execution time in cron format
	ReadConcurrency  uint   `json:"readConcurrency"` // Number of concurrent reads, defaults to CPU thread count
}

// Plan/Policy Status
const (
	RunningStatus = 1 // Status indicating the plan is running
	StopStatus    = 2 // Status indicating the plan is stopped
)
