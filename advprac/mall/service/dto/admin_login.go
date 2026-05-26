package dto

import (
	"fmt"
	"mall/utils/tools"
)

type GetVerifyCaptchaResp struct {
	Key             string `json:"key"`
	ImageBase64     string `json:"image_base64"`
	TileImageBase64 string `json:"tile_image_base_64"`
	TileHeight      int    `json:"tile_height"`
	TileWidth       int    `json:"tile_width"`
	TileX           int    `json:"tile_x"`
	TileY           int    `json:"tile_y"`
	Expire          int    `json:"expire"` // 过期时间
}

type GetVerifyCaptchaReq struct {
	Once string `url:"once"`
	Time int64  `url:"timestamp"`
	Sign string `url:"sign"` // 密钥加密: md5(once + zey2026+ ts)
}

type CheckCaptchaReq struct {
	Key    string `json:"key"`
	SlideX int    `json:"slide_x"`
	SlideY int    `json:"slide_y"`
}

type CheckCaptchaResp struct {
	Ticket string `json:"ticket"`
	Expire int64  `json:"expire"`
}

type MobilePasswordLoginReq struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Ticket   string `json:"ticket"`
}

type MobilePasswordLoginResp struct {
	Token string           `json:"token"`
	User  *GetUserInfoResp `json:"user"`
}

func (r *GetVerifyCaptchaReq) CheckSign() bool {
	return r.Sign == tools.Sha256Hash(fmt.Sprintf("%s%s%d", r.Once, "zey2026", r.Time))
}
