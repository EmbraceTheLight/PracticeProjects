package admin

import (
	"mall/adaptor"
	"mall/adaptor/repo/admin"
)

type Service struct {
	adminUSer admin.IAdminUser
}

func NewService(adaptor *adaptor.Adaptor) *Service {
	return &Service{
		adminUSer: admin.NewAdminUserRepo(adaptor),
	}
}
