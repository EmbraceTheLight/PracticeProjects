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
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.GetUserInfoReq{}
	err := ctx.BindQuery(req)
	if err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithError(err))
		return
	}
	repo, errno := c.adminSvc.GetUserInfo(ctx.Request.Context(), req)
	api.WriteResp(ctx, repo, errno)
}

func (c *Controller) CreateUser(ctx *gin.Context) {
	// 创建新的管理员的管理员信息
	adminUser := api.GetAdminUserFromCtx(ctx)
	if adminUser == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.CreateUserReq{}
	err := ctx.BindJSON(req)
	if err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithError(err))
		return
	}

	userID, errno := c.adminSvc.CreateUser(ctx.Request.Context(), adminUser, req)
	api.WriteResp(ctx, map[string]int64{
		"user_id": userID,
	}, errno)
}

func (c *Controller) UpdateUser(ctx *gin.Context) {
	// 更新管理员的管理员信息
	adminUser := api.GetAdminUserFromCtx(ctx)
	if adminUser == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.UpdateUserReq{}
	err := ctx.BindJSON(req)
	if err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithError(err))
		return
	}

	errno := c.adminSvc.UpdateUser(ctx.Request.Context(), adminUser, req)
	api.WriteResp(ctx, nil, errno)
}

func (c *Controller) UpdateUserStatus(ctx *gin.Context) {
	// 更新管理员的管理员信息
	adminUser := api.GetAdminUserFromCtx(ctx)
	if adminUser == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.UpdateUserStatusReq{}
	err := ctx.BindJSON(req)
	if err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithError(err))
		return
	}

	errno := c.adminSvc.UpdateUserStatus(ctx.Request.Context(), adminUser, req)
	api.WriteResp(ctx, nil, errno)
}
