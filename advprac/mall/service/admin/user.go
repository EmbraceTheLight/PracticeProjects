package admin

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mall/common"
	"mall/service/dto"
	"mall/utils/logger"
)

func (s *Service) GetUserInfo(ctx context.Context, adminUser *common.AdminUser) (*dto.GetUserInfoResp, common.Errno) {
	user, err := s.adminUSer.GetUserInfo(ctx, 1)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.UserNotFoundErr
		}
		logger.Error("GetUserInfo error", zap.Error(err), zap.Any("req", adminUser))
		return nil, common.DatabaseErr.WithError(err)
	}
	return &dto.GetUserInfoResp{
		UserID: user.ID,
		Name:   user.Name,
	}, common.Ok
}
