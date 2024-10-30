package api

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"time"
)

func pingHandler() iris.Handler {
	return func(ctx *context.Context) {
		ctx.Values().Set("data", time.Now())
	}
}

func AddPingRoute(app iris.Party) {
	// Used for health check
	app.Get("/ping", pingHandler())
}