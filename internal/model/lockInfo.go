package model

import "time"

type LockInfo struct {
	Time      time.Time `json:"time"`
	Exclusive bool      `json:"exclusive"`
	Hostname  string    `json:"hostname"`
	Username  string    `json:"username"`
	PID       int       `json:"pid"`
	UID       int       `json:"uid,omitempty"`
	GID       int       `json:"gid,omitempty"`
}
