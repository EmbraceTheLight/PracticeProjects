package verify

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"mall/adaptor"
	"mall/config"
	"mall/consts"
	"mall/service/do"
	auth "mall/utils/jwt"
)

type IVerify interface {
	SetCaptchaKey(ctx context.Context, key string, value string, expire time.Duration) error
	GetCaptchaKey(ctx context.Context, key string) (string, error)
	SetCaptchaTicket(ctx context.Context, key string, value string, expire time.Duration) error
	GetCaptchaTicket(ctx context.Context, key string) (string, error)

	GenerateToken(ctx context.Context, adminUser *do.GenerateTokenReq) (string, string, error)
	GetToken(ctx context.Context, userID int64, tokenType string) (string, error)
	GetAdminUserFromAccessToken(ctx context.Context, accessToken string) (*do.AdminUserTokenClaims, error)

	SaveToken(ctx context.Context, userID int64, token string, tokenType string, expire time.Duration) error
	IncreasePassWordErrCount(ctx context.Context, mobile string, expire time.Duration) (int64, error)
	DeletePassWordErrCount(ctx context.Context, mobile string) error
}

type verify struct {
	redis   *redis.Client
	jwtAuth *auth.JWTAuth
}

func NewVerify(adaptor adaptor.IAdaptor) IVerify {
	return &verify{
		redis:   adaptor.GetRedis(),
		jwtAuth: adaptor.GetJwtAuth(),
	}
}

func (v *verify) SetCaptchaKey(ctx context.Context, key string, value string, expire time.Duration) error {
	return v.redis.Set(ctx, fmtVerifyCaptchaKey(key), value, expire).Err()
}

func (v *verify) GetCaptchaKey(ctx context.Context, key string) (string, error) {
	redisKey := fmtVerifyCaptchaKey(key)
	captcha, err := v.redis.Get(ctx, redisKey).Result()
	if err != nil {
		return "", err
	}
	v.redis.Del(ctx, redisKey)
	return captcha, nil
}

func (v *verify) SetCaptchaTicket(ctx context.Context, key string, value string, expire time.Duration) error {
	return v.redis.Set(ctx, fmtVerifyCaptchaTicketKey(key), value, expire).Err()
}

func (v *verify) GetCaptchaTicket(ctx context.Context, key string) (string, error) {
	redisKey := fmtVerifyCaptchaTicketKey(key)
	ticket, err := v.redis.Get(ctx, redisKey).Result()
	if err != nil {
		return "", err
	}
	v.redis.Del(ctx, redisKey)
	return ticket, nil
}

func (v *verify) GenerateToken(ctx context.Context, adminUser *do.GenerateTokenReq) (string, string, error) {
	return v.jwtAuth.GenerateToken(ctx, adminUser)
}

func (v *verify) GetAdminUserFromAccessToken(ctx context.Context, accessToken string) (*do.AdminUserTokenClaims, error) {
	adminUser, err := v.parseAdminUserFromAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	storedToken, err := v.GetToken(ctx, adminUser.UserID, consts.AccessTokenKey)
	if err != nil {
		return nil, err
	}
	if storedToken != accessToken {
		return nil, fmt.Errorf("access token is not active")
	}
	return adminUser, nil
}

func (v *verify) parseAdminUserFromAccessToken(ctx context.Context, token string) (*do.AdminUserTokenClaims, error) {
	tokenClaim, err := v.jwtAuth.GetToken(ctx, token)
	if err != nil {
		return nil, err
	}
	adminUser, ok := tokenClaim.Claims.(*auth.AccessTokenClaims)
	if !ok {
		return nil, fmt.Errorf("token claim type error")
	}
	return &do.AdminUserTokenClaims{
		UserID:     adminUser.UserId,
		Name:       adminUser.Name,
		NickName:   adminUser.NickName,
		Sex:        adminUser.Sex,
		Status:     adminUser.Status,
		Mobile:     adminUser.Mobile,
		LarkOpenID: adminUser.LarkOpenID,
		UpdateAt:   adminUser.UpdateAt,
		CreateAt:   adminUser.CreateAt,
	}, nil
}

func (v *verify) IncreasePassWordErrCount(ctx context.Context, mobile string, expire time.Duration) (int64, error) {
	redisKey := fmtVerifyPasswordErrCount(mobile)
	pipeLine := v.redis.Pipeline()
	incr, err := pipeLine.Incr(ctx, redisKey).Result()
	if err != nil {
		return 0, err
	}
	if incr == 1 {
		pipeLine.Expire(ctx, redisKey, expire)
	}
	_, err = pipeLine.Exec(ctx)
	return incr, err
}

func (v *verify) DeletePassWordErrCount(ctx context.Context, mobile string) error {
	return v.redis.Del(ctx, fmtVerifyPasswordErrCount(mobile)).Err()
}

func (v *verify) SaveToken(ctx context.Context, userID int64, token string, tokenType string, expire time.Duration) error {
	switch tokenType {
	case consts.AccessTokenKey:
		return v.redis.Set(ctx, fmtUserAccessToken(userID), token, expire).Err()
	case consts.RefreshTokenKey:
		return v.redis.Set(ctx, fmtUserRefreshToken(userID), token, expire).Err()
	default:
		return fmt.Errorf("token type error")
	}
}

func (v *verify) GetToken(ctx context.Context, userID int64, tokenType string) (string, error) {
	switch tokenType {
	case consts.AccessTokenKey:
		return v.redis.Get(ctx, fmtUserAccessToken(userID)).Result()
	case consts.RefreshTokenKey:
		return v.redis.Get(ctx, fmtUserRefreshToken(userID)).Result()
	default:
		return "", fmt.Errorf("token type error")
	}
}

func fmtVerifyCaptchaKey(key string) string {
	return fmt.Sprintf("%s:captcha:%s", config.ServerFullName, key)
}

func fmtVerifyCaptchaTicketKey(key string) string {
	return fmt.Sprintf("%s:captcha:ticket:%s", config.ServerFullName, key)
}

func fmtUserAccessToken(key int64) string {
	return fmt.Sprintf("%s:access:token:%d", config.ServerFullName, key)
}

func fmtUserRefreshToken(key int64) string {
	return fmt.Sprintf("%s:refresh:token:%d", config.ServerFullName, key)
}

func fmtVerifyPasswordErrCount(mobile string) string {
	return fmt.Sprintf("%s:admin:user:password:errcount:%s", config.ServerFullName, mobile)
}
