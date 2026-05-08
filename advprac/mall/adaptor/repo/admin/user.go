package admin

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mall/adaptor"
	"mall/adaptor/repo/model"
	"mall/adaptor/repo/query"
)

type IAdminUser interface {
	GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error)
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

func (a *adminRepo) GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error) {
	qs := query.Use(a.db).AdminUser
	return qs.WithContext(ctx).Where(qs.ID.Eq(userId)).First()
}
