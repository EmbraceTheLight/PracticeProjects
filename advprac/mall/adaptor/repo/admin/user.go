package admin

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mall/adaptor"
	"mall/service/do"
)

type IAdminUser interface {
	HelloWorld(ctx context.Context, req *do.HelloWorld) (string, error)
}

type adminRepo struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewAdminUserRepo(adaptor *adaptor.Adaptor) IAdminUser {
	return &adminRepo{
		db:    adaptor.GetDB(),
		redis: adaptor.GetRedis(),
	}
}

func (a *adminRepo) HelloWorld(ctx context.Context, req *do.HelloWorld) (string, error) {
	return "hello world", nil
}
