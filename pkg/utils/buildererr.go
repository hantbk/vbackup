package utils

import (
	"fmt"

	"github.com/kataras/iris/context"

	"github.com/kataras/iris"
)

func ErrorCode(ctx *context.Context, code int, err error) {
	if err == nil {
		return
	}
	errstring := err.Error()
	ctx.StatusCode(code)
	ctx.Values().Set("message", errstring)
}

func Errore(ctx *context.Context, err error) {
	ErrorCode(ctx, iris.StatusBadRequest, err)
}

func ErrorStr(ctx *context.Context, err string) {
	Errore(ctx, fmt.Errorf(err))
}
