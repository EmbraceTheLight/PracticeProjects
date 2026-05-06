package api

import (
	"github.com/gin-gonic/gin"
	"mall/common"
	"net/http"
)

type Resp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	ErrMsg string `json:"err_msg"`
	Data   any    `json:"data"`
}

// WriteResp 统一的返回函数
func WriteResp(ctx *gin.Context, data any, errno common.Errno) {
	ctx.JSON(http.StatusOK, Resp{
		Code:   errno.Code,
		Msg:    errno.Msg,
		ErrMsg: errno.ErrMsg,
		Data:   data,
	})
}
