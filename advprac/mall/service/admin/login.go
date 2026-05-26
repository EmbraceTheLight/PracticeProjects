package admin

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"mall/common"
	"mall/service/dto"
	"mall/utils/logger"
)

func (s *Service) UserMobilePasswordLogin(ctx context.Context, req *dto.MobilePasswordLoginReq) (*dto.MobilePasswordLoginResp, common.Errno) {
	// 验证登录凭证, 该凭证仅在通过图形验证后才会出现
	_, err := s.verify.GetCaptchaTicket(ctx, req.Ticket)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, common.InvalidCaptchaErr
		}
		logger.Error("UserMobilePasswordLogin GetCaptchaTicket error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.RedisErr.WithError(err)
	}
	adminUser, err := s.adminUser.GetUserByMobile(ctx, req.Mobile)
	if err != nil {
		logger.Error("UserMobilePasswordLogin GetUserByMobile error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.DatabaseErr.WithError(err)
	}
	if adminUser == nil {
		return nil, common.UserNotFoundErr.WithMsg("手机号不存在, 请检查输入的手机号是否正确")
	}

	// 密码错误
	if adminUser.Password != req.Password {
		return nil, common.InvalidPasswordErr
	}

	// 密码正确, 返回
	return &dto.MobilePasswordLoginResp{
		Token: "token", // TODO: 使用 "token" 临时占位
		User: &dto.GetUserInfoResp{
			UserID: adminUser.ID,
			Name:   adminUser.Name,
		},
	}, common.Ok
}
