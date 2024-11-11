package utils

import (
	"time"

	"github.com/hantbk/vbackup/internal/model"
	"github.com/hantbk/vbackup/internal/server"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

var jwtSigner *jwt.Signer
var JwtKey string

func InitJwt() {
	JwtKey = server.Config().Jwt.Key
	t := server.Config().Jwt.MaxAge
	jwtMaxAge := time.Duration(t)
	jwtSigner = jwt.NewSigner(jwt.HS256, JwtKey, jwtMaxAge*time.Second)
}

func GetToken(data interface{}) (*model.TokenInfo, error) {
	token, err := jwtSigner.Sign(data)
	if err != nil {
		return nil, err
	}
	return &model.TokenInfo{
		Token:     string(token),
		ExpiresAt: time.Now().Add(jwtSigner.MaxAge),
	}, nil
}

func GetJwtVerifier() *jwt.Verifier {
	j := jwt.NewVerifier(jwt.HS256, JwtKey)
	j.ErrorHandler = func(ctx iris.Context, err error) {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{
			"error": err.Error(),
		})
	}
	return j
}
