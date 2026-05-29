package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mall/common"
	"mall/consts"
	"mall/utils/logger"
	"net/http"
)

type TokenFunc func(ctx context.Context, adminUserId string) (*common.User, error)
type TokenAdminFunc func(ctx context.Context, userId string) (*common.AdminUser, error)

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
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}
		// 用户鉴权逻辑
		user, err := getTokenFunc(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
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
		// 管理员鉴权逻辑
		adminUserId := ctx.GetHeader(consts.AdminTokenKey)
		if len(adminUserId) == 0 {
			logger.Error("AdminAuthMiddleware getAdminTokenFunc error", zap.String("msg", "token is empty"))
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}
		adminUser, err := getAdminTokenFunc(ctx, adminUserId)
		if err != nil {
			logger.Error("AdminAuthMiddleware getAdminTokenFunc error", zap.Error(err), zap.Int64("userId", adminUser.UserId), zap.String("username", adminUser.Name))
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}
		ctx.Set(consts.AdminUserKey, adminUser)
		ctx.Next()
	}
}
