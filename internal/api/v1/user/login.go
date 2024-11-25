package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/hantbk/vbackup/internal/consts"
	"github.com/hantbk/vbackup/internal/consts/globalcontext"
	"github.com/hantbk/vbackup/internal/entity/v1/oplog"
	"github.com/hantbk/vbackup/internal/entity/v1/sysuser"
	"github.com/hantbk/vbackup/internal/model"
	"github.com/hantbk/vbackup/internal/service/v1/common"
	logser "github.com/hantbk/vbackup/internal/service/v1/oplog"
	"github.com/hantbk/vbackup/internal/service/v1/user"
	"github.com/hantbk/vbackup/pkg/utils"
	"github.com/hantbk/vbackup/pkg/utils/otp"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

var userService user.Service
var logService logser.Service

var pwderrKey = "PwdErrCount:"

// Second
var lockTime = 1800

func init() {
	userService = user.GetService()
	logService = logser.GetService()
}

func loginHandler() iris.Handler {

	return func(ctx *context.Context) {
		var loginData model.LoginData
		err := ctx.ReadJSON(&loginData)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		udb := login(ctx, loginData.Username, loginData.Password)
		if udb == nil {
			return
		}

		udb.LastLogin = time.Now()
		userinfo := &model.Userinfo{
			Id:        udb.Id,
			Username:  udb.Username,
			NickName:  udb.NickName,
			Email:     udb.Email,
			Phone:     udb.Phone,
			Mfa:       false,
			LastLogin: udb.LastLogin.Format(consts.Custom),
		}

		if udb.OtpSecret != "" {
			userinfo.Mfa = true
		}
		// Set to 1 after closing the MFA
		if udb.OtpSecret == "1" {
			userinfo.Mfa = false
		}

		if userinfo.Mfa {
			if loginData.Code == "" {
				ctx.Values().Set("data", userinfo)
				return
			} else {
				if !otp.ValidCode(loginData.Code, udb.OtpInterval, udb.OtpSecret) {
					utils.Errore(ctx, err)
					return
				}
			}
		}
		token, err := utils.GetToken(userinfo)
		if err != nil {
			utils.ErrorStr(ctx, consts.TokenGenErrstr)
			return
		}
		userinfo.Token = token
		ctx.Values().Set("data", userinfo)

		// Store userinfo globally
		globalcontext.SetCurrentUser(userinfo)

		go func() {
			// Update last login log
			_ = userService.Update(udb, common.DBOptions{})
			var log oplog.OperationLog
			log.Operator = udb.Username
			log.Operation = "post"
			log.Url = "/api/v1/login"
			// Add new login log
			_ = logService.Create(&log, common.DBOptions{})
		}()
	}
}

func login(ctx *context.Context, username, password string) *sysuser.SysUser {
	u, err := userService.GetByUserName(username, common.DBOptions{})
	errk := pwderrKey + username
	count, ok := utils.Get(errk)
	if !ok || count == nil {
		count = 0
	}
	errcount := count.(int)
	if errcount >= 3 {
		utils.ErrorStr(ctx, fmt.Sprintf(consts.LockErrstr, errcount, lockTime/60))
		utils.Set(errk, errcount, lockTime)
		return nil
	}
	if err != nil || u == nil {
		utils.ErrorStr(ctx, consts.Pwderrstr)
		utils.Set(errk, errcount+1, lockTime)
		return nil
	}
	if !utils.ComparePwd(password, u.Password) {
		utils.ErrorStr(ctx, consts.Pwderrstr)
		utils.Set(errk, errcount+1, lockTime)
		return nil
	}

	ctx.Values().Set("userEmail", u.Email)

	return u
}

func refreshTokenHandler() iris.Handler {
	return func(ctx *context.Context) {
		curuser := utils.GetCurUser(ctx)
		token, err := utils.GetToken(curuser)
		if err != nil {
			utils.ErrorStr(ctx, "Token generation failed")
			return
		}
		ctx.Values().Set("data", token)
	}

}

func listHandler() iris.Handler {
	return func(ctx *context.Context) {
		var users []sysuser.SysUser
		users, err := userService.List(common.DBOptions{})
		if err != nil && err.Error() != "not found" {
			utils.Errore(ctx, err)
			return
		}
		for key, sysUser := range users {
			// Clear sensitive data
			sysUser.Password = ""
			sysUser.OtpSecret = ""
			sysUser.OtpInterval = 0
			users[key] = sysUser
		}
		ctx.Values().Set("data", users)
	}
}

func updateHandler() iris.Handler {
	return func(ctx *context.Context) {
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		var muser model.Userinfo
		err = ctx.ReadJSON(&muser)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}

		sysUser, err := userService.Get(id, common.DBOptions{})
		if err != nil {
			return
		}
		sysUser.NickName = muser.NickName
		sysUser.Email = muser.Email
		sysUser.Phone = muser.Phone
		err = userService.Update(sysUser, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", "")
	}
}

func createHandler() iris.Handler {
	return func(ctx *context.Context) {
		var muser sysuser.SysUser
		err := ctx.ReadJSON(&muser)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		if strings.TrimSpace(muser.Password) == "" {
			utils.ErrorStr(ctx, "Password cannot be empty")
			return
		}
		if len(muser.Password) < 6 {
			utils.ErrorStr(ctx, "Password length cannot be less than 6 characters")
			return
		}
		encodePWD, err := utils.EncodePWD(muser.Password)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		muser.Password = encodePWD
		err = userService.Create(&muser, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", "")
	}
}

func delHandler() iris.Handler {
	return func(ctx *context.Context) {
		id, err := ctx.Params().GetInt("id")
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		sysUser, err := userService.Get(id, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		if sysUser.Username == "admin" {
			utils.ErrorStr(ctx, "The admin account cannot be deleted!")
			return
		}
		curu := utils.GetCurUser(ctx)
		if curu.Id == sysUser.Id {
			utils.ErrorStr(ctx, "You cannot delete yourself!")
			return
		}
		err = userService.DeleteStruct(sysUser, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", "")
	}
}

func repwdHandler() iris.Handler {
	return func(ctx *context.Context) {
		var pwdData model.RePwdData
		err := ctx.ReadJSON(&pwdData)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		if strings.TrimSpace(pwdData.Password) == "" {
			utils.ErrorStr(ctx, "Password cannot be empty")
			return
		}
		if len(pwdData.Password) < 6 {
			utils.ErrorStr(ctx, "Password length must be at least 6 characters")
			return
		}
		curuser := utils.GetCurUser(ctx)
		sysUser, err := userService.Get(curuser.Id, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		if !utils.ComparePwd(pwdData.OldPassword, sysUser.Password) {
			utils.ErrorStr(ctx, consts.RePwderrstr)
			return
		}
		newPwd, err := utils.EncodePWD(pwdData.Password)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		sysUser.Password = newPwd
		err = userService.Update(sysUser, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", "Modification successful")
	}
}

func otpHandler() iris.Handler {
	return func(ctx *context.Context) {
		curu := utils.GetCurUser(ctx)
		getOtp, err := otp.GetOtp(curu.Username, "vbackup", 30)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", getOtp)
	}
}

func bindOtpHandler() iris.Handler {
	return func(ctx *context.Context) {
		var otpInfo model.OtpInfo
		err := ctx.ReadJSON(&otpInfo)
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		res := otp.ValidCode(otpInfo.Code, otpInfo.Interval, otpInfo.Secret)
		if !res {
			utils.ErrorStr(ctx, "Verification code is incorrect!")
			return
		}
		curu := utils.GetCurUser(ctx)
		sysUser, err := userService.Get(curu.Id, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		sysUser.OtpSecret = otpInfo.Secret
		sysUser.OtpInterval = otpInfo.Interval
		err = userService.Update(sysUser, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", "Binding successful")
	}
}

func putOtpHandler() iris.Handler {
	return func(ctx *context.Context) {
		curu := utils.GetCurUser(ctx)
		err := userService.ClearOtp(curu.Username, common.DBOptions{})
		if err != nil {
			utils.Errore(ctx, err)
			return
		}
		ctx.Values().Set("data", "Shutdown successful")
	}
}

func Install(parent iris.Party) {
	// User-related APIs
	sp := parent.Party("/")
	// Login
	sp.Post("/login", loginHandler())
	// Refresh token
	sp.Post("/refreshToken", refreshTokenHandler())
	// Get user list
	sp.Get("/user", listHandler())
	// Delete user
	sp.Delete("/user/:id", delHandler())
	// Update user
	sp.Put("/user/:id", updateHandler())
	// Create user
	sp.Post("/user", createHandler())
	// Change password
	sp.Post("/repwd", repwdHandler())
	// Get OTP QR code
	sp.Get("/otp", otpHandler())
	// Bind OTP
	sp.Post("/otp", bindOtpHandler())
	// Update OTP
	sp.Put("/otp", putOtpHandler())
}
