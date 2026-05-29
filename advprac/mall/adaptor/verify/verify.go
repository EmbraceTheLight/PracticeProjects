package verify

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"mall/adaptor"
	"mall/config"
	"mall/consts"
	"mall/service/do"
	auth "mall/utils/jwt"
	"strconv"
	"time"
)

type IVerify interface {
	SetCaptchaKey(ctx context.Context, key string, value string, expire time.Duration) error
	GetCaptchaKey(ctx context.Context, key string) (string, error)
	SetCaptchaTicket(ctx context.Context, key string, value string, expire time.Duration) error
	GetCaptchaTicket(ctx context.Context, key string) (string, error)

	// GenerateToken 根据管理员用户信息生成 access token 与 refresh token
	GenerateToken(ctx context.Context, adminUser *do.GenerateTokenReq) (string, string, error)

	// GetToken 从 redis 中获取 token. 根据 tokenType 决定获取 refresh-token 还是 access-token
	GetToken(ctx context.Context, adminUserId string, tokenType string) (string, error)

	// GetAdminUserFromToken 从 token 解析出管理员用户信息
	GetAdminUserFromToken(ctx context.Context, adminUserId string) (*do.GetAdminUserFromTokenResp, error)

	SaveToken(ctx context.Context, userId int64, token string, tokenType string, expire time.Duration) error
	IncreasePassWordErrCount(ctx context.Context, mobile string, expire time.Duration) (int64, error)
	DeletePassWordErrCount(ctx context.Context, mobile string) error
}

type verify struct {
	redis   *redis.Client
	jwtAuth *auth.JWTAuth
}

func (v *verify) SetCaptchaKey(ctx context.Context, key string, value string, expire time.Duration) error {
	redisKey := fmtVerifyCaptchaKey(key)
	return v.redis.Set(ctx, redisKey, value, expire).Err()
}

func (v *verify) GetCaptchaKey(ctx context.Context, key string) (string, error) {
	redisKey := fmtVerifyCaptchaKey(key)
	captcha, err := v.redis.Get(ctx, redisKey).Result()
	if err != nil {
		return "", err
	}

	// 每个 captcha 只允许用一次, 从 redis 中获取完成后, 就要将其删掉
	v.redis.Del(ctx, redisKey)
	return captcha, nil
}

func (v *verify) SetCaptchaTicket(ctx context.Context, key string, value string, expire time.Duration) error {
	redisKey := fmtVerifyCaptchaTicketKey(key)
	return v.redis.Set(ctx, redisKey, value, expire).Err()
}

func (v *verify) GetCaptchaTicket(ctx context.Context, key string) (string, error) {
	redisKey := fmtVerifyCaptchaTicketKey(key)
	ticket, err := v.redis.Get(ctx, redisKey).Result()
	if err != nil {
		return "", err
	}

	// 每个 captcha ticket 只允许用一次, 和 captcha 一样
	v.redis.Del(ctx, redisKey)
	return ticket, nil
}

func (v *verify) GetAdminUserFromToken(ctx context.Context, adminUserId string) (*do.GetAdminUserFromTokenResp, error) {
	accessToken, err := v.GetToken(ctx, adminUserId, consts.AccessTokenKey)
	if err != nil {
		return nil, err
	}
	return v.getAdminUserFromToken(ctx, accessToken)
}

func (v *verify) getAdminUserFromToken(ctx context.Context, token string) (*do.GetAdminUserFromTokenResp, error) {
	tokenClaim, err := v.jwtAuth.GetToken(ctx, token)
	if err != nil {
		return nil, err
	}
	adminUser, ok := tokenClaim.Claims.(*auth.AccessTokenClaims)
	if !ok {
		return nil, fmt.Errorf("token claim type error")
	}
	return &do.GetAdminUserFromTokenResp{
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

func (v *verify) GenerateToken(ctx context.Context, adminUser *do.GenerateTokenReq) (string, string, error) {
	return v.jwtAuth.GenerateToken(ctx, adminUser)
}
func (v *verify) DeletePassWordErrCount(ctx context.Context, mobile string) error {
	redisKey := fmtVerifyPasswordErrCount(mobile)
	return v.redis.Del(ctx, redisKey).Err()
}

func (v *verify) SaveToken(ctx context.Context, userId int64, token string, tokenType string, expire time.Duration) error {
	switch tokenType {
	case consts.AccessTokenKey:
		redisKey := fmtUserAccessToken(userId)
		return v.redis.Set(ctx, redisKey, token, expire).Err()
	case consts.RefreshTokenKey:
		redisKey := fmtUserRefreshToken(userId)
		return v.redis.Set(ctx, redisKey, token, expire).Err()
	default:
		return fmt.Errorf("token type error")
	}
}

func (v *verify) GetToken(ctx context.Context, adminUserId string, tokenType string) (string, error) {
	id, _ := strconv.ParseInt(adminUserId, 10, 64)
	switch tokenType {
	case consts.AccessTokenKey:
		redisKey := fmtUserAccessToken(id)
		return v.redis.Get(ctx, redisKey).Result()
	case consts.RefreshTokenKey:
		redisKey := fmtUserRefreshToken(id)
		return v.redis.Get(ctx, redisKey).Result()
	default:
		return "", fmt.Errorf("token type error")
	}
}

func NewVerify(adaptor adaptor.IAdaptor) IVerify {
	return &verify{
		redis:   adaptor.GetRedis(),
		jwtAuth: adaptor.GetJwtAuth(),
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
