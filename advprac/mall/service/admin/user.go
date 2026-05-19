package admin

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mall/common"
	"mall/service/do"
	"mall/service/dto"
	"mall/utils/logger"
)

func (s *Service) GetUserInfo(ctx context.Context, req *dto.GetUserInfoReq) (*dto.GetUserInfoResp, common.Errno) {
	user, err := s.adminUSer.GetUserInfo(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.UserNotFoundErr
		}
		logger.Error("GetUserInfo error", zap.Error(err), zap.Any("req", req))
		return nil, common.DatabaseErr.WithError(err)
	}
	return &dto.GetUserInfoResp{
		UserID: user.ID,
		Name:   user.Name,
	}, common.Ok
}

func (s *Service) CreateUser(ctx context.Context, adminUser *common.AdminUser, req *dto.CreateUserReq) (int64, common.Errno) {
	userID, err := s.adminUSer.CreateUser(ctx, &do.CreateUser{
		CreateBy: adminUser.UserId,
		Name:     req.Name,
		NickName: req.NickName,
		Mobile:   req.Mobile,
		Sex:      req.Sex,
	})
	if err != nil {
		logger.Error("CreateUser error", zap.Error(err), zap.Any("req", req))
		return -1, common.DatabaseErr.WithError(err)
	}
	return userID, common.Ok
}

func (s *Service) UpdateUser(ctx context.Context, adminUser *common.AdminUser, req *dto.UpdateUserReq) common.Errno {
	err := s.adminUSer.UpdateUser(ctx, &do.UpdateUser{
		UpdateBy: adminUser.UserId,
		ID:       req.ID,
		Name:     req.Name,
		NickName: req.NickName,
		Sex:      req.Sex,
	})
	if err != nil {
		logger.Error("UpdateUser error", zap.Error(err), zap.Any("req", req))
		return common.DatabaseErr.WithError(err)
	}
	return common.Ok
}

func (s *Service) UpdateUserStatus(ctx context.Context, adminUser *common.AdminUser, req *dto.UpdateUserStatusReq) common.Errno {
	err := s.adminUSer.UpdateUserStatus(ctx, &do.UpdateUserStatus{
		UpdateBy: adminUser.UserId,
		ID:       req.ID,
		Status:   req.Status,
	})
	if err != nil {
		logger.Error("UpdateUserStatus error", zap.Error(err), zap.Any("req", req))
		return common.DatabaseErr.WithError(err)
	}
	return common.Ok
}
