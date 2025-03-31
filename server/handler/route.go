package handler

import (
	_interface "github.com/duanchi/min/v2/interface"
	"github.com/duanchi/min/v2/server/httpserver"
	"github.com/duanchi/min/v2/server/httpserver/context"
	"github.com/duanchi/min/v2/server/middleware"
	"reflect"
)

func RouteHandle(path string, handle reflect.Value, ctx *context.Context, engine *httpserver.Httpserver) {
	handle.Interface().(_interface.RouterInterface).Handle(ctx.Request().Path(), ctx.Request().Method(), ctx.Params(), ctx)

	handlers := middleware.GetHandlersBeforeResponse()

	for _, handler := range handlers {
		handler(ctx)
	}
}
