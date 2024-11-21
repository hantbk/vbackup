package utils

import (
	"time"

	"github.com/hantbk/vbackup/internal/model"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/jwt"
)

func GetCurUser(ctx *context.Context) *model.Userinfo {
	ctx.User()
	var u *model.Userinfo
	if ctx.GetHeader("Authorization") != "" {
		u = jwt.Get(ctx).(*model.Userinfo)
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
