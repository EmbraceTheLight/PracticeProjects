package common

type AdminUser struct {
	UserId int64  `json:"user_id"`
	Name   string `json:"name"`
}

type User struct {
	UserId   int64  `json:"user_id"`
	NickName string `json:"nick_ame"`
}
