package abstract

import (
	_interface "github.com/duanchi/min/v2/interface"
	"github.com/duanchi/min/v2/server/httpserver/context"
	"github.com/duanchi/min/v2/types/middleware"
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

func (this *Middleware) BeforeRoute(ctx *context.Context) {
	ctx.Next()
}

func (this *Middleware) AfterRoute(ctx *context.Context) {
	ctx.Next()
}

func (this *Middleware) BeforeResponse(ctx *context.Context) {
	ctx.Next()
}

func (this *Middleware) AfterResponse(ctx *context.Context) {
	ctx.Next()
}

func (this *Middleware) AfterPanic(ctx *context.Context) {
	ctx.Next()
}
