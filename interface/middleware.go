package _interface

import (
	"github.com/duanchi/min/v2/server/httpserver/context"
	"github.com/duanchi/min/v2/types/middleware"
)

type MiddlewareInterface interface {
	Includes() (includes middleware.Includes)
	Excludes() (excludes middleware.Excludes)
	BeforeRoute(ctx *context.Context)
	AfterRoute(ctx *context.Context)
	BeforeResponse(ctx *context.Context)
	AfterResponse(ctx *context.Context)
	AfterPanic(ctx *context.Context)
}
