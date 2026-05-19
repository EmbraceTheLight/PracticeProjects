package router

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"mall/consts"
	"mall/utils/logger"
	"time"
)

// GetRequestBody 获取请求体
func GetRequestBody(ctx *gin.Context) string {
	data, _ := io.ReadAll(ctx.Request.Body)
	return string(data)
}

type responseWriterWrapper struct {
	// 这里定义匿名 gin.ResponseWriter 是为了满足 gin.ResponseWriter 定义的其他接口
	gin.ResponseWriter
	multiWriter io.Writer
}

func (r *responseWriterWrapper) Write(b []byte) (int, error) {
	return r.multiWriter.Write(b)
}

// GetResponseBody 获取到响应数据, 如果响应长度超过 1024 字节, 则进行截取, 只取前 1024 字节的响应
func GetResponseBody(respBuffer *bytes.Buffer) string {
	rawResp := respBuffer.String()
	if len(rawResp) > 1024 {
		rawResp = rawResp[:1024]
	}
	return rawResp
}

func AccessLogMiddleware(filter func(c *gin.Context) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 没有定义过滤函数, 或过滤函数结果为 false, 则该接口不记录日志
		if filter != nil && filter(c) == false {
			c.Next()
			return
		}
		begin := time.Now()
		body := GetRequestBody(c)
		c.Request.Body = io.NopCloser(bytes.NewBufferString(body)) // 重新赋值请求体 Body, 上一步读取请求体 body 会导致其被清空
		fields := []zap.Field{
			zap.String("ip", c.RemoteIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("params", c.Request.URL.RawQuery),
			zap.String("body", body),
			zap.String("token", c.GetHeader(consts.UserTokenKey)),
		}

		var responseBuffer bytes.Buffer

		// io.MultiWriter 将响应数据定向到两个位置:
		// c.Writer: gin 正常写入响应数据
		// responseBuffer: 自定义字段, 接收响应数据
		multiWriter := io.MultiWriter(c.Writer, &responseBuffer)
		c.Writer = &responseWriterWrapper{
			c.Writer,
			multiWriter,
		}
		c.Next()

		// 记录请求消耗的时间, 单位: 毫秒
		elapsedTimeStr := fmt.Sprintf("%d ms", time.Since(begin).Milliseconds())
		fields = append(fields, zap.String("duration", elapsedTimeStr))

		// 记录响应状态码
		fields = append(fields, zap.Int("status", c.Writer.Status()))

		// 记录响应数据
		fields = append(fields, zap.String("resp", GetResponseBody(&responseBuffer)))
		logger.Info("access log", fields...)
	}
}
