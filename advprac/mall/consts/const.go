package consts

import "time"

const (
	UserTokenKey       = "token"
	CustomerUserKey    = "user_key"
	AdminUserKey       = "admin_user_key"
	AccessTokenKey     = "access_token"
	AccessTokenExpire  = 5 * time.Hour
	RefreshTokenKey    = "refresh_token"
	RefreshTokenExpire = 30 * 24 * time.Hour
)
const (
	AdminUserTokenExpire   = 24 * time.Hour
	PasswordErrCountExpire = 5 * time.Minute
)

const (
	PasswordErrCountMaxLimit = 3
)

const (
	Enable  = 1
	Disable = -1
)
