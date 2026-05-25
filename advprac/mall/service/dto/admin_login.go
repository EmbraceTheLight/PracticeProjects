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
	Expire          int    `json:"expire"`
}

type GetVerifyCaptchaReq struct {
	Once string `url:"once"`
	Time int64  `url:"timestamp"`
	Sign string `url:"sign"` // 密钥加密: md5(once + zey2026+ ts)
}

func (r *GetVerifyCaptchaReq) CheckSign() bool {
	return r.Sign == tools.Sha256Hash(fmt.Sprintf("%s%s%d", r.Once, "zey2026", r.Time))
}
