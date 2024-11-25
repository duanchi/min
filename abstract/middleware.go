package abstract

import (
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/duanchi/min/types/middleware"
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

func (this *Middleware) BeforeRoute(ctx *context.Context) {}

func (this *Middleware) AfterRoute(ctx *context.Context) {}

func (this *Middleware) BeforeResponse(ctx *context.Context) {}

func (this *Middleware) AfterResponse(ctx *context.Context) {}

func (this *Middleware) AfterPanic(ctx *context.Context) {}
