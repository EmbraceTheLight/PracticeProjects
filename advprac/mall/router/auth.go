package router

import (
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(filter func(c *gin.Context) bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
