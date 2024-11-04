package utils

import (
	"time"

	"github.com/hantbk/vbackup/internal/model"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/jwt"
)

func GetCurUser(ctx *context.Context) *model.UserInfo {
	ctx.User()
	var u *model.UserInfo
	if ctx.GetHeader("Authorization") != "" {
		u = jwt.Get(ctx).(*model.UserInfo)
	}
	return u
}

func GetTokenExpires(ctx *context.Context) time.Time {
	if ctx.GetHeader("Authorization") != "" {
		vt := jwt.GetVerifiedToken(ctx)
		if vt != nil {
			return vt.StandardClaims.ExpiresAt()
		}
	}
	return time.Now()
}
