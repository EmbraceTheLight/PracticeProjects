package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"mall/adaptor"
	"mall/config"
	"time"
)

type IVerify interface {
	SetCaptchaKey(ctx context.Context, key string, value string, expire time.Duration) error
	GetCaptchaKey(ctx context.Context, key string) (string, error)
	SetCaptchaTicket(ctx context.Context, key string, value string, expire time.Duration) error
	GetCaptchaTicket(ctx context.Context, key string) (string, error)
}

type verify struct {
	redis *redis.Client
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

func NewVerify(adaptor adaptor.IAdaptor) IVerify {
	return &verify{
		redis: adaptor.GetRedis(),
	}
}

func fmtVerifyCaptchaKey(key string) string {
	return fmt.Sprintf("%s:captcha:%s", config.ServerFullName, key)
}
func fmtVerifyCaptchaTicketKey(key string) string {
	return fmt.Sprintf("%s:captcha:ticket:%s", config.ServerFullName, key)
}
