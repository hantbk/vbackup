package operation

import (
	"github.com/hantbk/vbackup/internal/entity/v1/common"
	wsTaskInfo "github.com/hantbk/vbackup/internal/store/ws_task_info"
)

type Operation struct {
	common.BaseModel `storm:"inline"`
	RepositoryId     int                  `json:"repositoryId"`
	PolicyId         int                  `json:"policyId"`
	Type             int                  `json:"type"`
	Status           int                  `json:"status"`
	Logs             []*wsTaskInfo.Sprint `json:"logs"`
}

const (
	CHECK_TYPE        = 1 // CHECK - Check repository status
	REBUILDINDEX_TYPE = 2 // REBUILDINDEX - Rebuild index
	PRUNE_TYPE        = 3 // PRUNE - Clean up unused data
	FORGET_TYPE       = 4 // FORGET - Remove expired snapshots
	MIGRATE_TYPE      = 5 // MIGRATE - Migrate
)
