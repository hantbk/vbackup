package ws

import (
	"github.com/hantbk/vbackup/internal/store/log"
	info "github.com/hantbk/vbackup/internal/store/task"
	task "github.com/hantbk/vbackup/internal/store/ws_task_info"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func Install(parent iris.Party) {
	// WebSocket endpoint
	wsParty := parent.Party("/ws")
	// Backup task SockJS endpoint
	taskh := task.CreateTaskHandler("/task/sockjs", info.TaskInfos)
	wsParty.Any("/task/sockjs/{p:path}", func(context *context.Context) {
		taskh.ServeHTTP(context.ResponseWriter(), context.Request())
	})

	logh := task.CreateTaskHandler("/log/sockjs", log.LogInfos)
	wsParty.Any("/log/sockjs/{p:path}", func(context *context.Context) {
		logh.ServeHTTP(context.ResponseWriter(), context.Request())
	})
}
