package admin

import (
	"context"
	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"mall/common"
	"mall/service/dto"
	"mall/utils/logger"
	"mall/utils/tools"
	"time"
)

// GetSlideCaptcha 获取图形验证码数据
func (s *Service) GetSlideCaptcha(ctx context.Context) (*dto.GetVerifyCaptchaResp, common.Errno) {
	captchaData, err := s.captcha.Generate()
	if err != nil {
		logger.Error("GetSlideCaptcha Generate error", zap.Error(err))
		return nil, common.ServeErr.WithError(err)
	}
	dotData := captchaData.GetData()
	if dotData == nil {
		logger.Error("GetSlideCaptcha Generate error", zap.Error(err))
		return nil, common.ServeErr.WithMsg("GetData is nil")
	}

	dots, err := sonic.Marshal(dotData)
	if err != nil {
		logger.Error("GetSlideCaptcha Generate error", zap.Error(err))
		return nil, common.ServeErr.WithError(err)
	}

	var mainB64Data, tileB64Data string // mainB64Data: 滑块图片(主图) base64编码, tileB64Data: 滑块背景图(缩略图) base64编码
	mainB64Data, err = captchaData.GetMasterImage().ToBase64()
	if err != nil {
		logger.Error("GetSlideCaptcha Generate error", zap.Error(err))
		return nil, common.ServeErr.WithError(err)
	}

	tileB64Data, err = captchaData.GetTileImage().ToBase64()
	if err != nil {
		logger.Error("GetSlideCaptcha Generate error", zap.Error(err))
		return nil, common.ServeErr.WithError(err)
	}

	key := tools.UUIDHex()
	err = s.verify.SetCaptchaKey(ctx, key, string(dots), 2*time.Minute)
	if err != nil {
		logger.Error("GetSlideCaptcha SetCaptchaKey error", zap.Error(err))
		return nil, common.RedisErr.WithError(err)
	}
	return &dto.GetVerifyCaptchaResp{
		Key:             key,
		ImageBase64:     mainB64Data,
		TileImageBase64: tileB64Data,
		TileHeight:      dotData.Height,
		TileWidth:       dotData.Width,
		TileX:           dotData.DX,
		TileY:           dotData.DY,
		Expire:          110,
	}, common.Ok
}
