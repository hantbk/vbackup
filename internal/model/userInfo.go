package model

import "time"

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"` // OTP code
}

type RePwdData struct {
	OldPassword string `json:"oldPassword"`
	Password    string `json:"password"`
}

type UserInfo struct {
	Id        int        `json:"id"`
	Username  string     `json:"username"`
	Nickname  string     `json:"nickname"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	LastLogin string     `json:"lastLogin"`
	Mfa       bool       `json:"mfa"`
	Token     *TokenInfo `json:"token"`
}

type TokenInfo struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
