package wsTaskInfo

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hantbk/vbackup/internal/server"
	"github.com/hantbk/vbackup/pkg/utils"
	"github.com/kataras/iris/v12"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type WsTaskInfo interface {
	GetId() int
	SetId(id int)
	SetBound(c chan string)
	GetBound() chan string
	CloseBound()
	IntoBound(msg string)
	SetSockJSSession(sockjs.Session)
	SendMsg(msg interface{})
	CloseSockJSSession(reason string, status uint32)
}

type Message struct {
	Id int
}

type WsTask interface {
	Get(id int) WsTaskInfo
	Set(id int, task WsTaskInfo)
	Close(id int, reason string, status uint32)
	GetCount() int
}

func CreateTaskHandler(path string, wsTask WsTask) http.Handler {
	return sockjs.NewHandler(path, sockjs.DefaultOptions, taskHandler(wsTask))
}

func taskHandler(wsTask WsTask) func(session sockjs.Session) {
	return func(session sockjs.Session) {
		var (
			buf      string
			err      error
			msg      Message
			taskInfo WsTaskInfo
		)
		if buf, err = session.Recv(); err != nil {
			server.Logger().Errorf("taskHandler: can't Recv: %v", err)
			return
		}
		if err = json.Unmarshal([]byte(buf), &msg); err != nil {
			server.Logger().Errorf("taskHandler: can't UnMarshal (%v): %s", err, buf)
			return
		}
		id := msg.Id
		if taskInfo = wsTask.Get(id); taskInfo == nil {
			errstr := utils.ToJSONString(iris.Map{
				"success": false,
				"code":    iris.StatusBadRequest,
				"message": "There are no ongoing tasks for this ID",
			})
			err := session.Send(errstr)
			if err != nil {
				return
			}
			return
		}
		taskInfo.SetSockJSSession(session)
		wsTask.Set(id, taskInfo)
		t := time.Now().Second()
		taskInfo.IntoBound(string(rune(id + t)))
	}
}
