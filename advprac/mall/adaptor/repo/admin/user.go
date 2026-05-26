package admin

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mall/adaptor"
	"mall/adaptor/repo/model"
	"mall/adaptor/repo/query"
	"mall/consts"
	"mall/service/do"
)

type IAdminUser interface {
	GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error)
	GetUserByMobile(ctx context.Context, mobile string) (*model.AdminUser, error)

	CreateUser(ctx context.Context, req *do.CreateUser) (int64, error)
	UpdateUser(ctx context.Context, req *do.UpdateUser) error
	UpdateUserStatus(ctx context.Context, req *do.UpdateUserStatus) error
}

type adminRepo struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewAdminUserRepo(adaptor adaptor.IAdaptor) IAdminUser {
	return &adminRepo{
		db:    adaptor.GetDB(),
		redis: adaptor.GetRedis(),
	}
}

func (a *adminRepo) GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error) {
	qs := query.Use(a.db).AdminUser
	return qs.WithContext(ctx).Where(qs.ID.Eq(userId)).First()
}

func (a *adminRepo) CreateUser(ctx context.Context, req *do.CreateUser) (int64, error) {
	qs := query.Use(a.db).AdminUser
	user := &model.AdminUser{
		Name:     req.Name,
		NickName: req.NickName,
		Mobile:   req.Mobile,
		Sex:      req.Sex,
		Status:   consts.Enable, // 启用管理员
		CreateBy: req.CreateBy,
	}
	err := qs.WithContext(ctx).Create(user)
	if err != nil {
		return -1, err
	}
	return user.ID, nil
}

func (a *adminRepo) UpdateUser(ctx context.Context, req *do.UpdateUser) error {
	qs := query.Use(a.db).AdminUser
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Updates(&model.AdminUser{
		Name:     req.Name,
		NickName: req.NickName,
		Sex:      req.Sex,
		UpdateBy: req.UpdateBy,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *adminRepo) UpdateUserStatus(ctx context.Context, req *do.UpdateUserStatus) error {
	qs := query.Use(a.db).AdminUser
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Updates(&model.AdminUser{
		Status:   req.Status,
		UpdateBy: req.UpdateBy,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *adminRepo) GetUserByMobile(ctx context.Context, mobile string) (*model.AdminUser, error) {
	qs := query.Use(a.db).AdminUser
	return qs.WithContext(ctx).Where(qs.Mobile.Eq(mobile)).First()
}
