package admin

import (
	"github.com/gin-gonic/gin"
	"mall/adaptor"
	"mall/api"
	"mall/common"
	"mall/service/admin"
	"mall/service/dto"
)

type Controller struct {
	adaptor *adaptor.Adaptor

	// services
	adminSvc *admin.Service
}

func NewCtrl(adaptor *adaptor.Adaptor) *Controller {
	return &Controller{
		adaptor:  adaptor,
		adminSvc: admin.NewService(adaptor),
	}
}

func (c *Controller) HelloWorld(ctx *gin.Context) {
	resp, errno := c.adminSvc.HelloWorld(ctx.Request.Context(), &common.AdminUser{}, &dto.HelloWorldReq{Name: "ad"})
	api.WriteResp(ctx, resp, errno)
}
