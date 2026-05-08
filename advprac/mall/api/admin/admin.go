package admin

import (
	"github.com/gin-gonic/gin"
	"mall/adaptor"
	"mall/api"
	"mall/common"
	"mall/service/admin"
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

func (c *Controller) GetUserInfo(ctx *gin.Context) {
	// 获取 token
	repo, errno := c.adminSvc.GetUserInfo(ctx.Request.Context(), &common.AdminUser{})
	api.WriteResp(ctx, repo, errno)
}
