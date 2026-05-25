package admin

import (
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
