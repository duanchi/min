package abstract

import (
	_interface "github.com/duanchi/min/interface"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Bean
	_interface.MiddlewareInterface
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
