package _interface

import (
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/duanchi/min/types/middleware"
)

type MiddlewareInterface interface {
	Includes() (includes middleware.Includes)
	Excludes() (excludes middleware.Excludes)
	BeforeRoute(ctx *context.Context) error
	AfterRoute(ctx *context.Context) error
	BeforeResponse(ctx *context.Context) error
	AfterResponse(ctx *context.Context) error
	AfterPanic(ctx *context.Context) error
}
