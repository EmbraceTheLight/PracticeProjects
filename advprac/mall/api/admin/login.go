package admin

import (
	"context"
	"github.com/gin-gonic/gin"
	"mall/api"
	"mall/common"
	"mall/service/dto"
)

func (c *Controller) GetSmsCodeCaptcha(ctx *gin.Context) {
	req := &dto.GetVerifyCaptchaReq{}
	if err := ctx.BindQuery(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithError(err))
		return
	}
	captchaData, errno := c.adminSvc.GetSlideCaptcha(ctx.Request.Context())
	api.WriteResp(ctx, captchaData, errno)
}

func (c *Controller) CheckSmsCodeCaptcha(ctx *gin.Context) {
	req := &dto.CheckCaptchaReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithError(err))
		return
	}
	captchaData, errno := c.adminSvc.CheckSlideCaptcha(ctx.Request.Context(), req)
	api.WriteResp(ctx, captchaData, errno)
}

func (c *Controller) MobilePasswordLogin(ctx *gin.Context) {
	req := &dto.MobilePasswordLoginReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithError(err))
		return
	}
	userData, errno := c.adminSvc.UserMobilePasswordLogin(ctx.Request.Context(), &dto.MobilePasswordLoginReq{
		Mobile:   req.Mobile,
		Password: req.Password,
		Ticket:   req.Ticket,
	})
	api.WriteResp(ctx, userData, errno)
}

func (c *Controller) GetAdminUserByToken(ctx context.Context, token string) (*common.AdminUser, common.Errno) {
	adminUser, err := c.adminSvc.GetAdminUserByToken(ctx, token)
	if err.NotOk() {
		return nil, err
	}
	return adminUser, common.Ok
}
