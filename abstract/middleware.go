package abstract

import (
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/types/middleware"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Bean
	_interface.MiddlewareInterface
}

func (this *Middleware) Includes() (includes middleware.Includes) {
	return
}
func (this *Middleware) Excludes() (excludes middleware.Excludes) {
	return
}

func (this *Middleware) BeforeRoute(ctx *gin.Context) {
	ctx.Next()
}

func (this *Middleware) AfterRoute(ctx *gin.Context) {
	ctx.Next()
}

func (this *Middleware) BeforeResponse(ctx *gin.Context) {
	ctx.Next()
}

func (this *Middleware) AfterResponse(ctx *gin.Context) {
	ctx.Next()
}

func (this *Middleware) AfterPanic(ctx *gin.Context) {
	ctx.Next()
}
