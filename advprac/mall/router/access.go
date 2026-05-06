package router

import "github.com/gin-gonic/gin"

func AccessLogMiddleware(filter func(c *gin.Context) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 没有定义过滤函数, 或过滤函数结果为 false, 则该接口不记录日志
		if filter != nil && filter(c) == false {
			c.Next()
			return
		}

		// 记录日志
		c.Next()
	}
}
