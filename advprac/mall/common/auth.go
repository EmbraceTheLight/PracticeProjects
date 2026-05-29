package common

type AdminUser struct {
	UserId     int64  `json:"user_id"`
	Name       string `json:"name"`
	NickName   string `json:"nick_name"`
	Sex        int64  `json:"sex"`
	Status     int64  `json:"status"`
	Mobile     string `json:"mobile"`
	LarkOpenID string `json:"lark_open_id"`
}

type User struct {
	UserId   int64  `json:"user_id"`
	NickName string `json:"nick_ame"`
}
