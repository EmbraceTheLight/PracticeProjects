package router

// AdminAuthWhiteList 所有中间件通用的路由白名单
var AdminAuthWhiteList map[string]bool = map[string]bool{
	"/ping":                                  true,
	"/metrics":                               true,
	"/api/mall/admin/v1/user/verify/captcha": true,
}
