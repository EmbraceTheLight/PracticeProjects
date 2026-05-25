package admin

import (
	"github.com/wenlng/go-captcha/v2/slide"
	"mall/adaptor"
	"mall/adaptor/redis"
	"mall/adaptor/repo/admin"
	"mall/utils/captcha"
)

type Service struct {
	adminUSer admin.IAdminUser
	verify    redis.IVerify
	captcha   slide.Captcha
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		adminUSer: admin.NewAdminUserRepo(adaptor),
		verify:    redis.NewVerify(adaptor),
		captcha:   captcha.NewSlideCaptcha(),
	}
}
