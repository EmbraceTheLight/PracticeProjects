package router

import (
	"github.com/gin-gonic/gin"
	"net/http/pprof"
)

func SetupPprof(r *gin.Engine, prefix string) {
	group := r.Group(prefix)
	{
		group.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		group.GET("/profile", gin.WrapF(pprof.Profile))
		group.GET("/symbol", gin.WrapF(pprof.Symbol))
		group.GET("/trace", gin.WrapF(pprof.Trace))
		group.GET("/allocs", pprofHandler("allocs"))
		group.GET("/block", pprofHandler("block"))
		group.GET("/goroutine", pprofHandler("goroutine"))
		group.GET("/heap", pprofHandler("heap"))
		group.GET("/mutex", pprofHandler("mutex"))
	}
}

func pprofHandler(name string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler(name).ServeHTTP(ctx.Writer, ctx.Request)
	}
}
