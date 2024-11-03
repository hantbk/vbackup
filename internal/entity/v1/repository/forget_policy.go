package repository

import "github.com/hantbk/vbackup/internal/entity/v1/common"

type ForgetPolicy struct {
	common.BaseModel `storm:"inline"`
	RepositoryId     int    `json:"repositoryId"` // ID of the associated repository
	Path             string `json:"path"`         // Path to apply the forget policy
	Status           int    `json:"status"`       // Status of the policy

	/**
	Forget policy type:
	- "last"    : Keep the last few snapshots
	- "hourly"  : Retain snapshots by hour
	- "daily"   : Retain snapshots by day
	- "weekly"  : Retain snapshots by week
	- "monthly" : Retain snapshots by month
	- "yearly"  : Retain snapshots by year
	*/
	Type string `json:"type"` // Type of retention policy
	Value int `json:"value"` // Retention value
}

