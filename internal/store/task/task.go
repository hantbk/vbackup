package task

import (
	wsTaskInfo "github.com/hantbk/vbackup/internal/store/ws_task_info"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

// Task Status
const (
	StatusNew     = 0 // New
	StatusRunning = 1 // Running
	StatusEnd     = 2 // Completed
	StatusError   = 3 // Error
)

var TaskInfos = &TaskMap{TaskInfos: make(map[int]wsTaskInfo.WsTaskInfo)}

type TaskInfo struct {
	id            int
	bound         chan string
	sockJSSession sockjs.Session
}
