package task

import (
	"github.com/hantbk/vbackup/internal/entity/v1/common"
	"github.com/hantbk/vbackup/internal/model"
)

type Task struct {
	common.BaseModel `storm:"inline"`
	Name             string               `json:"name"`          // Task name
	Path             string               `json:"path"`          // Backup or restore path
	PlanId           int                  `json:"planId"`        // ID of the associated plan
	RepositoryId     int                  `json:"repositoryId"`  // ID of the associated repository
	Status           int                  `json:"status"`        // Task status
	ParentId         string               `json:"parentId"`      // Parent snapshot ID, if any
	Scanner          *model.VerboseUpdate `json:"scanner"`       // Scan results
	ScannerError     *model.ErrorUpdate   `json:"scannerError"`  // Scan errors
	ArchivalError    []model.ErrorUpdate  `json:"archivalError"` // List of errors encountered during backup
	Summary          *model.SummaryOutput `json:"summary"`       // Summary of backup operation
	Progress         *model.StatusUpdate  `json:"progress"`      // Current progress update
	RestoreError     []model.ErrorUpdate  `json:"restoreError"`  // List of errors encountered during restore
	ReadConcurrency  uint                 // Number of concurrent reads, default is 2
}
