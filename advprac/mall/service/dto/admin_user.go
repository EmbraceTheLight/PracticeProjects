package dto

import "time"

type UserInfoDto struct {
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

type GetUserInfoReq struct {
	ID int64 `form:"id"`
}

type CreateUserReq struct {
	Name     string `json:"name"`
	NickName string `json:"nick_name"`
	Mobile   string `json:"mobile"`
	Sex      int64  `json:"sex"`
}

type UpdateUserReq struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	NickName string `json:"nick_name"`
	Sex      int64  `json:"sex"`
}

type UpdateUserStatusReq struct {
	ID     int64 `json:"id"`
	Status int64 `json:"status"`
}
