package router

import (
	"context"
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
	checkFunc func() error
	admin     *admin.Controller
	customer  *customer.Controller
}

func NewRouter(adaptor adaptor.IAdaptor, conf *config.Config, checkFunc func() error) *Router {
	return &Router{
		FullPPROF: conf.HttpServer.EnableFullPPROF,
		rootPath:  "/api/mall",
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
	return true
}

func (r *Router) route(root *gin.RouterGroup) {
	r.adminRoute(root)
	r.customerRoute(root)
}

func (r *Router) customerRoute(root *gin.RouterGroup) {
	cstRoot := root.Group("/customer", AuthMiddleware(r.SpanFilter, func(ctx context.Context, token string) (*common.User, error) {
		return nil, nil
	}))
	cstRoot.GET("/user/info", r.admin.GetUserInfo)

}

func (r *Router) adminRoute(root *gin.RouterGroup) {
	adminRoot := root.Group("/admin", AdminAuthMiddleware(r.SpanFilter, func(ctx context.Context, token string) (*common.AdminUser, error) {
		return nil, nil
	}))
	adminRoot.GET("/user/info", r.admin.GetUserInfo)
}
