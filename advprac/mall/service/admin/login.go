package admin

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"mall/common"
	"mall/consts"
	"mall/service/do"
	"mall/service/dto"
	"mall/utils/logger"
	"strconv"
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

	// 账户未启用
	if adminUser.Status != consts.Enable {
		return nil, common.InvalidPasswordErr.WithMsg("用户未启用, 请联系管理员")
	}

	// 密码错误次数限制. 5 分钟内同一个用户不能输错 3 次密码
	if adminUser.Password != req.Password {
		errCount, err := s.verify.IncreasePassWordErrCount(ctx, req.Mobile, consts.PasswordErrCountExpire)
		if err != nil {
			logger.Error("UserMobilePasswordLogin IncreasePassWordErrCount error", zap.Error(err), zap.String("mobile", req.Mobile))
			return nil, common.RedisErr.WithError(err)
		}
		if errCount >= consts.PasswordErrCountMaxLimit {
			return nil, common.PasswordErrLimit
		}
	}

	// 密码正确, 删除密码错误计数
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

	// 生成 token 并存入 redis 中, 之后前端只需携带 user_id 作为请求头, 后端即可通过查询 token 进行验证
	err = s.processToken(ctx, adminUserDto)
	if err != nil {
		logger.Error("UserMobilePasswordLogin processToken error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.JwtCreateTokenErr.WithError(err)
	}
	// 密码正确, 返回.
	return &dto.MobilePasswordLoginResp{
		User: adminUserDto,
	}, common.Ok
}

func (s *Service) GetAdminUser(ctx context.Context, adminUserId string) (*common.AdminUser, error) {
	adminUser, err := s.verify.GetAdminUserFromToken(ctx, adminUserId)
	if err != nil {
		// token 过期, 查看 refresh-token 是否过期, 如果没有过期, 则重新生成 access-token 并正常返回
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, redis.Nil) {
			// 尝试重置 access-token
			err = s.refreshToken(ctx, adminUserId)
			if err != nil {
				logger.Error("GetAdminUser::refreshToken error", zap.Error(err), zap.String("adminUserId", adminUserId))
				return nil, err
			}

			// 根据新生成的 access-token, 重新获取用户信息
			adminUser, err = s.verify.GetAdminUserFromToken(ctx, adminUserId)
			if err != nil {
				logger.Error("GetAdminUser::GetAdminUserFromToken error", zap.Error(err), zap.String("adminUserId", adminUserId))
				return nil, common.AuthErr.WithError(errors.New("请重新登录"))
			}
		} else {
			logger.Error("GetAdminUser error", zap.Error(err), zap.String("adminUserId", adminUserId))
			return nil, common.AuthErr.WithError(errors.New("请重新登录"))
		}
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

func (s *Service) processToken(ctx context.Context, adminUser *dto.UserInfoDto) error {
	// 暂不使用 refresh-token, 只生成并保存 access-token
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
		return err
	}
	err = s.verify.SaveToken(ctx, adminUser.UserID, accessToken, consts.AccessTokenKey, consts.AccessTokenExpire)
	if err != nil {
		return err
	}

	err = s.verify.SaveToken(ctx, adminUser.UserID, refreshToken, consts.RefreshTokenKey, consts.RefreshTokenExpire)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) refreshToken(ctx context.Context, adminUserId string) error {
	_, err := s.verify.GetToken(ctx, adminUserId, consts.RefreshTokenKey)
	if err != nil {
		return err
	}

	// 重新获取用户信息, 以生成 access-token
	id, _ := strconv.ParseInt(adminUserId, 10, 64)
	adminUserInfo, err := s.adminUser.GetUserInfo(ctx, id)
	if err != nil {
		logger.Error("refreshToken::GetUserInfo error", zap.Error(err), zap.String("adminUserId", adminUserId))
		return err
	}

	// 生成并保存 access-token
	accessToken, _, err := s.verify.GenerateToken(ctx, &do.GenerateTokenReq{
		UserID:     adminUserInfo.ID,
		Name:       adminUserInfo.Name,
		NickName:   adminUserInfo.NickName,
		Sex:        adminUserInfo.Sex,
		Status:     adminUserInfo.Status,
		Mobile:     adminUserInfo.Mobile,
		LarkOpenID: adminUserInfo.LarkOpenID,
		UpdateAt:   adminUserInfo.UpdateAt,
		CreateAt:   adminUserInfo.CreateAt,
	})
	if err != nil {
		logger.Error("refreshToken::GenerateToken error", zap.Error(err), zap.String("adminUserId", adminUserId))
		return err
	}

	err = s.verify.SaveToken(ctx, adminUserInfo.ID, accessToken, consts.AccessTokenKey, consts.AccessTokenExpire)
	if err != nil {
		logger.Error("refreshToken::SaveToken error", zap.Error(err), zap.String("adminUserId", adminUserId))
		return err
	}
	return nil
}
