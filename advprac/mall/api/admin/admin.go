package admin

import (
	"github.com/gin-gonic/gin"
	"mall/adaptor"
	"mall/api"
	"mall/common"
	"mall/service/admin"
)

type Controller struct {
	adaptor adaptor.IAdaptor

	// services
	adminSvc *admin.Service
}

func NewCtrl(adaptor adaptor.IAdaptor) *Controller {
	return &Controller{
		adaptor:  adaptor,
		adminSvc: admin.NewService(adaptor),
	}
}

func (c *Controller) GetUserInfo(ctx *gin.Context) {
	// 获取用户信息, 二次校验. 理论上中间件会拦截用户不存在的情况
	// 这里做二次校验, 逻辑上更加严谨
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthError)
		return
	}
	repo, errno := c.adminSvc.GetUserInfo(ctx.Request.Context(), &common.AdminUser{})
	api.WriteResp(ctx, repo, errno)
}
