package sysuser

import (
	"time"

	"github.com/hantbk/vbackup/internal/entity/v1/common"
)

type SysUser struct {
	common.BaseModel `storm:"inline"`
	Username         string    `json:"username" storm:"index,unique"` // Username
	NickName         string    `json:"nickName"`                      // Nickname
	Password         string    `json:"password"`                      // Password
	Email            string    `json:"email"`                         // Email
	Phone            string    `json:"phone"`                         // Phone
	OtpSecret        string    `json:"OtpSecret"`                     // OTP secret
	OtpInterval      int       `json:"otpInterval"`                   // OTP interval
	LastLogin        time.Time `json:"lastLogin"`                     // Last login time
}
