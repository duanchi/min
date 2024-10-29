package handler

import (
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/server/httpserver"
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/duanchi/min/server/middleware"
	"reflect"
)

func RouteHandle(path string, handle reflect.Value, ctx *context.Context, engine *httpserver.Httpserver) {
	handle.Interface().(_interface.RouterInterface).Handle(ctx.Request().Path(), ctx.Request().Method(), ctx.Params(), ctx)

	handlers := middleware.GetHandlersBeforeResponse()

	for _, handler := range handlers {
		handler(ctx)
	}
}
