package router

import (
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 普通用户鉴权
func AuthMiddleware(filter func(c *gin.Context) bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// AdminAuthMiddleware 管理员用户鉴权
func AdminAuthMiddleware(filter func(c *gin.Context) bool) gin.HandlerFunc {
	// 返回一个匿名函数作为中间件处理函数
	return func(ctx *gin.Context) {

		// 中间件函数体目前为空，实际使用时应该包含认证逻辑
	}
}
