package oplog

import "github.com/hantbk/vbackup/internal/entity/v1/common"

type OperationLog struct {
	common.BaseModel `storm:"inline"`
	Operator         string `json:"operator"`  // Operator (the person performing the action)
	Operation        string `json:"operation"` // Action performed
	Url              string `json:"url"`       // Request URL
	Data             string `json:"data"`      // Request data
}
