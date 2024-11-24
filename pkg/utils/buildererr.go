package utils

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// ErrorCode sets the HTTP status code and a message from the error.
func ErrorCode(ctx *context.Context, code int, err error) {
	if err == nil {
		return
	}

	errstring := err.Error()
	// Log the error (optional)
	ctx.Application().Logger().Errorf("Error occurred: %s", errstring)

	// Set HTTP status and response
	ctx.StatusCode(code)
	ctx.JSON(map[string]string{
		"message": errstring,
	})
}

// Errore wraps an error and responds with a Bad Request status code.
func Errore(ctx *context.Context, err error) {
	ErrorCode(ctx, iris.StatusBadRequest, err)
}

// ErrorStr wraps a string as an error and passes it to Errore.
func ErrorStr(ctx *context.Context, err string) {
	Errore(ctx, fmt.Errorf(err))
}
