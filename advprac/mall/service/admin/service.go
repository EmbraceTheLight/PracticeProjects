package admin

import (
	"github.com/wenlng/go-captcha/v2/slide"
	"mall/adaptor"
	"mall/adaptor/repo/admin"
	"mall/adaptor/verify"
	"mall/utils/captcha"
)

type Service struct {
	adminUser admin.IAdminUser
	verify    verify.IVerify
	captcha   slide.Captcha
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		adminUser: admin.NewAdminUserRepo(adaptor),
		verify:    verify.NewVerify(adaptor),
		captcha:   captcha.NewSlideCaptcha(),
	}
}
