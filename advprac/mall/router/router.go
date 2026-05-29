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
	if AdminAuthWhiteList[ctx.Request.URL.Path] == true {
		return true
	}
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
	cstRoot := root.Group("/customer", AuthMiddleware(r.SpanFilter, func(ctx context.Context, userId string) (*common.User, error) {
		return nil, nil
	}))
	cstRoot.GET("/user/info", r.admin.GetUserInfo)

}

func (r *Router) adminRoute(root *gin.RouterGroup) {
	adminRoot := root.Group("/admin", AdminAuthMiddleware(r.SpanFilter, func(ctx context.Context, adminUserId string) (*common.AdminUser, error) {
		adminUser, err := r.admin.GetUserAdminByToken(ctx, adminUserId)
		return adminUser, err
	}))
	// 登录, 不应鉴权, 添加到白名单
	adminRoot.GET("/v1/user/verify/captcha", r.admin.GetSmsCodeCaptcha)
	adminRoot.POST("/v1/user/verify/captcha/check", r.admin.CheckSmsCodeCaptcha)
	adminRoot.POST("/v1/user/mobile/password_login", r.admin.MobilePasswordLogin)
	adminRoot.GET("/v1/user/info", r.admin.GetUserInfo)
	adminRoot.POST("/v1/user/create", r.admin.CreateUser)
	adminRoot.POST("/v1/user/update", r.admin.UpdateUser)
	adminRoot.POST("/v1/user/status/update", r.admin.UpdateUserStatus)
}
