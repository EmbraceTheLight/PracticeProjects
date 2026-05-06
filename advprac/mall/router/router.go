package router

import (
	"github.com/gin-gonic/gin"
	"mall/adaptor"
	"mall/api"
	"mall/api/admin"
	"mall/api/customer"
	"mall/common"
	"mall/config"
	"net/http"
)

type IRouter interface {
	Register(engine *gin.Engine)
	SpanFilter(ctx *gin.Context) bool
	AccessRecordFilter(ctx *gin.Context) bool
}

type Router struct {
	FullPPROF bool
	rootPath  string
	conf      *config.Config
	adaptor   *adaptor.Adaptor
	checkFunc func() error
	admin     *admin.Controller
	customer  *customer.Controller
}

func NewRouter(adaptor *adaptor.Adaptor, conf *config.Config, checkFunc func() error) *Router {
	return &Router{
		FullPPROF: conf.HttpServer.EnableFullPPROF,
		rootPath:  "/api/mall",
		adaptor:   adaptor,
		conf:      conf,
		checkFunc: checkFunc,
		admin:     admin.NewCtrl(adaptor),
		customer:  customer.NewCtrl(adaptor),
	}
}

func (r *Router) checkServer() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		err := r.checkFunc()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		api.WriteResp(ctx, nil, common.Ok)
	}
}

func (r *Router) Register(app *gin.Engine) {
	app.Use(AuthMiddleware(r.SpanFilter))
	if r.conf.HttpServer.EnableFullPPROF {
		SetupPprof(app, "/debug/pprof")
	}
	app.Any("/ping", r.checkServer())
	root := app.Group(r.rootPath)
	r.route(root)
}

func (r *Router) SpanFilter(ctx *gin.Context) bool {
	return false
}
func (r *Router) AccessRecordFilter(ctx *gin.Context) bool {
	return false
}

func (r *Router) route(root *gin.RouterGroup) {
	root.GET("/hello", r.admin.HelloWorld)
}
