package wsTaskInfo

import (
	"io"
	"sync"
	"time"
)

// "MaxErrorNum" is the maximum amount of data stored, used in conjunction with "limitNum."
const MaxErrorNum = 10

type Sprintf struct {
	Sprints        []*Sprint
	taskInfo       WsTaskInfo
	MinUpdatePause time.Duration
	lastUpdate     time.Time
	limitNum       int
	limitNumLock   sync.Mutex
	io.Writer
}

type Sprint struct {
	Clear bool   `json:"clear"`
	Text  string `json:"text"`
	Time  string `json:"time"`
	Level int    `json:"level"`
}

const (
	Info    = 1
	Warning = 2
	Success = 3
	Error   = 4
)
