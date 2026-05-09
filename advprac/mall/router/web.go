package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"mall/utils/logger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type App struct {
	server *gin.Engine
	addr   string
}

func NewApp(httpPort int, router IRouter) *App {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Recovery 中间件, 捕获全局 panic
	engine.Use(gin.Recovery())

	// 日志中间件, 使用自定义过滤器, 因为某些接口不需要记录日志
	engine.Use(AccessLogMiddleware(router.AccessRecordFilter))

	// 注册业务路由
	router.Register(engine)
	return &App{
		server: engine,
		addr:   ":" + strconv.Itoa(httpPort),
	}
}

func (app *App) Run() {
	srv := &http.Server{
		Addr:    app.addr,
		Handler: app.server,
	}

	// 异步启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Listen err: %v", err)
		}
	}()

	closeCh := make(chan os.Signal, 1)
	signal.Notify(closeCh, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	msg := <-closeCh
	logger.Warn("server closing: ", zap.String("msg", msg.String()))
	srv.Shutdown(context.TODO())
}
