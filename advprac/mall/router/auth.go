package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mall/common"
	"mall/consts"
	auth "mall/utils/jwt"
	"mall/utils/logger"
	"net/http"
	"strings"
)

type TokenFunc func(ctx context.Context, adminUserId string) (*common.User, error)
type TokenAdminFunc func(ctx context.Context, accessToken string) (*common.AdminUser, error)

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
		accessToken := getBearerToken(ctx)
		if len(accessToken) == 0 {
			logger.Error("AdminAuthMiddleware getAdminTokenFunc error", zap.String("msg", "token is empty"))
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}
		adminUser, err := getAdminTokenFunc(ctx, accessToken)
		if err != nil {
			logger.Error("AdminAuthMiddleware getAdminTokenFunc error", zap.Error(err))
			ctx.JSON(http.StatusUnauthorized, common.AuthErr)
			ctx.Abort()
			return
		}
		ctx.Set(consts.AdminUserKey, adminUser)
		ctx.Next()
	}
}

func getBearerToken(ctx *gin.Context) string {
	authHeader := strings.TrimSpace(ctx.GetHeader(auth.AccessTokenHeader))
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], auth.BearerTokenPrefix) {
		return ""
	}
	return parts[1]
}
