package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"mall/common"
	"mall/consts"
	"net/http"
)

type TokenFunc func(ctx context.Context, token string) (*common.User, error)
type TokenAdminFunc func(ctx context.Context, token string) (*common.AdminUser, error)

// AuthMiddleware 普通用户鉴权
func AuthMiddleware(filter func(c *gin.Context) bool, getTokenFunc TokenFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 接口位于白名单中
		if filter != nil && filter(ctx) != false {
			ctx.Next()
			return
		}
		token := ctx.GetHeader(consts.UserTokenKey)
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, common.AuthError)
			ctx.Abort()
			return
		}
		// 用户鉴权逻辑
		user, err := getTokenFunc(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.AuthError)
			ctx.Abort()
			return
		}
		ctx.Set(consts.CustomerUserKey, user)
		ctx.Next()
	}
}

// AdminAuthMiddleware 管理员用户鉴权
func AdminAuthMiddleware(filter func(c *gin.Context) bool, getAdminTokenFunc TokenAdminFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 接口位于白名单中
		if filter != nil && filter(ctx) != false {
			ctx.Next()
			return
		}
		token := ctx.GetHeader(consts.AdminTokenKey)
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, common.AuthError)
			ctx.Abort()
			return
		}
		// 管理员鉴权逻辑
		adminUser, err := getAdminTokenFunc(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.AuthError)
			ctx.Abort()
			return
		}
		ctx.Set(consts.AdminUserKey, adminUser)
		ctx.Next()
	}
}
