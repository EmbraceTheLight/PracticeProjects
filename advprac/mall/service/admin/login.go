package admin

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"mall/common"
	"mall/consts"
	"mall/service/do"
	"mall/service/dto"
	"mall/utils/logger"
)

func (s *Service) UserMobilePasswordLogin(ctx context.Context, req *dto.MobilePasswordLoginReq) (*dto.MobilePasswordLoginResp, common.Errno) {
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
		return nil, common.UserNotFoundErr.WithMsg("mobile not found")
	}

	if adminUser.Status != consts.Enable {
		return nil, common.InvalidPasswordErr.WithMsg("user is disabled")
	}

	if adminUser.Password != req.Password {
		errCount, err := s.verify.IncreasePassWordErrCount(ctx, req.Mobile, consts.PasswordErrCountExpire)
		if err != nil {
			logger.Error("UserMobilePasswordLogin IncreasePassWordErrCount error", zap.Error(err), zap.String("mobile", req.Mobile))
			return nil, common.RedisErr.WithError(err)
		}
		if errCount >= consts.PasswordErrCountMaxLimit {
			return nil, common.PasswordErrLimit
		}
		return nil, common.InvalidPasswordErr
	}

	_ = s.verify.DeletePassWordErrCount(ctx, req.Mobile)
	adminUserDto := &dto.UserInfoDto{
		UserID:     adminUser.ID,
		Name:       adminUser.Name,
		NickName:   adminUser.NickName,
		Sex:        adminUser.Sex,
		Status:     adminUser.Status,
		Mobile:     adminUser.Mobile,
		LarkOpenID: adminUser.LarkOpenID,
		UpdateAt:   adminUser.UpdateAt,
		CreateAt:   adminUser.CreateAt,
	}

	accessToken, refreshToken, err := s.processToken(ctx, adminUserDto)
	if err != nil {
		logger.Error("UserMobilePasswordLogin processToken error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.JwtCreateTokenErr.WithError(err)
	}

	return &dto.MobilePasswordLoginResp{
		User:         adminUserDto,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, common.Ok
}

func (s *Service) GetAdminUserByAccessToken(ctx context.Context, accessToken string) (*common.AdminUser, error) {
	adminUser, err := s.verify.GetAdminUserFromAccessToken(ctx, accessToken)
	if err != nil {
		logger.Error("GetAdminUserByAccessToken error", zap.Error(err))
		return nil, common.AuthErr.WithError(errors.New("please login again"))
	}
	return &common.AdminUser{
		UserId:     adminUser.UserID,
		Name:       adminUser.Name,
		NickName:   adminUser.NickName,
		Sex:        adminUser.Sex,
		Status:     adminUser.Status,
		Mobile:     adminUser.Mobile,
		LarkOpenID: adminUser.LarkOpenID,
	}, nil
}

func (s *Service) processToken(ctx context.Context, adminUser *dto.UserInfoDto) (string, string, error) {
	accessToken, refreshToken, err := s.verify.GenerateToken(ctx, &do.GenerateTokenReq{
		UserID:     adminUser.UserID,
		Name:       adminUser.Name,
		NickName:   adminUser.NickName,
		Sex:        adminUser.Sex,
		Status:     adminUser.Status,
		Mobile:     adminUser.Mobile,
		LarkOpenID: adminUser.LarkOpenID,
		UpdateAt:   adminUser.UpdateAt,
		CreateAt:   adminUser.CreateAt,
	})
	if err != nil {
		return "", "", err
	}

	err = s.verify.SaveToken(ctx, adminUser.UserID, accessToken, consts.AccessTokenKey, consts.AccessTokenExpire)
	if err != nil {
		return "", "", err
	}

	err = s.verify.SaveToken(ctx, adminUser.UserID, refreshToken, consts.RefreshTokenKey, consts.RefreshTokenExpire)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
