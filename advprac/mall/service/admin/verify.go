package admin

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/wenlng/go-captcha/v2/slide"
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

func (s *Service) CheckSlideCaptcha(ctx context.Context, req *dto.CheckCaptchaReq) (*dto.CheckCaptchaResp, common.Errno) {
	captData, err := s.verify.GetCaptchaKey(ctx, req.Key)
	if err != nil {
		logger.Error("CheckSlideCaptcha GetCaptchaKey error", zap.Error(err))
		return nil, common.RedisErr.WithError(err)
	}
	if captData == "" {
		return nil, common.ParamErr.WithMsg("滑块已过期, 请刷新重试")
	}
	dot := slide.Block{}
	err = sonic.Unmarshal([]byte(captData), &dot)
	if err != nil {
		logger.Error("CheckSlideCaptcha Unmarshal error", zap.Error(err))
		return nil, common.InvalidCaptchaErr.WithError(err)
	}
	ok := slide.Validate(req.SlideX, req.SlideY, dot.X, dot.Y, 5)
	if !ok {
		return nil, common.InvalidCaptchaErr
	}

	// 图形验证通过, 设置 ticket, 作为后续登录票据凭证, 防止恶意刷登录接口
	// 前端在后续发送登录请求时携带该 ticket
	ticket := tools.UUIDHex()
	err = s.verify.SetCaptchaTicket(ctx, ticket, req.Key, 6*time.Minute)
	if err != nil {
		logger.Error("CheckSlideCaptcha SetCaptchaTicket error", zap.Error(err))
		return nil, common.RedisErr.WithError(err)
	}
	return &dto.CheckCaptchaResp{Ticket: ticket, Expire: 320}, common.Ok
}
