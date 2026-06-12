package do

import "time"

type CreateUser struct {
	CreateBy int64  `json:"create_by"`
	Name     string `json:"name"`
	NickName string `json:"nick_name"`
	Mobile   string `json:"mobile"`
	Sex      int64  `json:"sex"`
}

type UpdateUser struct {
	UpdateBy int64  `json:"update_by"`
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	NickName string `json:"nick_name"`
	Sex      int64  `json:"sex"`
}

type UpdateUserStatus struct {
	UpdateBy int64 `json:"update_by"`
	ID       int64 `json:"id"`
	Status   int64 `json:"status"`
}

type GenerateTokenReq struct {
	UserID     int64     `json:"user_id"`
	Name       string    `json:"name"`
	NickName   string    `json:"nick_name"`
	Sex        int64     `json:"sex"`
	Status     int64     `json:"status"`
	Mobile     string    `json:"mobile"`
	LarkOpenID string    `json:"lark_open_id"`
	UpdateAt   time.Time `json:"update_at"`
	CreateAt   time.Time `json:"create_at"`
}

type AdminUserTokenClaims struct {
	UserID     int64     `json:"user_id"`
	Name       string    `json:"name"`
	NickName   string    `json:"nick_name"`
	Sex        int64     `json:"sex"`
	Status     int64     `json:"status"`
	Mobile     string    `json:"mobile"`
	LarkOpenID string    `json:"lark_open_id"`
	UpdateAt   time.Time `json:"update_at"`
	CreateAt   time.Time `json:"create_at"`
}
