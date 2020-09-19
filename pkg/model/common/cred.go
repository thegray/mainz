package common

import "time"

type UserInfo struct {
	UserID     int
	UserName   string
	UUIDAccess string
}

type AuthInfo struct {
	UserAgent string
	IP        string
	AuthTime  time.Time
	ExpTime   int64
}
