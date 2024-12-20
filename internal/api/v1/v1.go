package v1

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/kataras/iris/v12/context"

	"github.com/hantbk/vbackup/internal/api/v1/dashboard"
	"github.com/hantbk/vbackup/internal/api/v1/operation"
	"github.com/hantbk/vbackup/internal/api/v1/plan"
	"github.com/hantbk/vbackup/internal/api/v1/policy"
	"github.com/hantbk/vbackup/internal/api/v1/repository"
	"github.com/hantbk/vbackup/internal/api/v1/restic"
	"github.com/hantbk/vbackup/internal/api/v1/system"
	"github.com/hantbk/vbackup/internal/api/v1/task"
	"github.com/hantbk/vbackup/internal/api/v1/user"
	"github.com/hantbk/vbackup/internal/api/v1/ws"
	"github.com/hantbk/vbackup/internal/consts/system_status"
	"github.com/hantbk/vbackup/internal/entity/v1/oplog"
	"github.com/hantbk/vbackup/internal/model"
	"github.com/hantbk/vbackup/internal/server"
	"github.com/hantbk/vbackup/internal/service/v1/common"
	logser "github.com/hantbk/vbackup/internal/service/v1/oplog"
	"github.com/hantbk/vbackup/pkg/utils"
	"github.com/kataras/iris/v12"
)

// Authentication API Whitelist.
var apiWhiteList = WhiteList{"post:/api/v1/login", "get:/api/v1/system/version", "get:/api/v1/system/version/latest"}

type WhiteList []string

func (w WhiteList) In(name string) bool {
	for i := range w {
		if w[i] == name {
			return true
		}
	}
	return false
}

// jwtHandler authenticates the user with JWT.
func jwtHandler() iris.Handler {
	verifier := utils.GetJwtVerifier()
	middleware := verifier.Verify(func() interface{} { return new(model.Userinfo) })
	return func(c *context.Context) {
		path := c.Path()
		method := strings.ToLower(c.Method())
		if apiWhiteList.In(method + ":" + path) {
			c.Next()
			return
		}
		middleware(c)
	}
}

func logHandler() iris.Handler {
	return func(ctx *context.Context) {
		method := strings.ToLower(ctx.Method())
		if method != "post" && method != "delete" && method != "put" {
			ctx.Next()
			return
		}
		path := ctx.Path()
		if path == "/api/v1/refreshToken" {
			ctx.Next()
			return
		}
		if apiWhiteList.In(method + ":" + path) {
			ctx.Next()
			return
		}
		curuser := utils.GetCurUser(ctx)
		if curuser == nil {
			resp := iris.Map{
				"success":      true,
				"systemStatus": system_status.Normal,
				"data":         ctx.Values().Get("data"),
				"isDocker":     server.IsDocker(),
			}

			ctx.StatusCode(iris.StatusUnauthorized)
			_ = ctx.JSON(resp, iris.JSON{})
			return
		}

		var log oplog.OperationLog
		log.Operator = curuser.Username
		log.Operation = method
		log.Url = path
		if method == "post" || method == "put" {
			data, _ := ctx.GetBody()
			log.Data = string(data)
			if path == "/api/v1/repwd" {
				log.Data = ""
			}
			if path == "/api/v1/user" {
				log.Data = ""
			}
			ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(data))
		}
		logService := logser.GetService()
		go func() {
			_ = logService.Create(&log, common.DBOptions{})
		}()
		ctx.Next()
	}
}

func AddV1Route(app iris.Party) {
	// API v1 collection
	v1Party := app.Party("/v1")
	v1Party.Use(jwtHandler())
	v1Party.Use(logHandler())
	dashboard.Install(v1Party)
	user.Install(v1Party)
	restic.Install(v1Party)
	system.Install(v1Party)
	repository.Install(v1Party)
	plan.Install(v1Party)
	task.Install(v1Party)
	operation.Install(v1Party)
	policy.Install(v1Party)
	ws.Install(v1Party)
}
