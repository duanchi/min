package handler

import (
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/server/middleware"
	"github.com/gin-gonic/gin"
	"reflect"
)

func RouteHandle(path string, handle reflect.Value, ctx *gin.Context, engine *gin.Engine) {
	params := ctx.Params
	method := ctx.Request.Method

	handle.Interface().(_interface.RouterInterface).Handle(ctx.Request.URL.Path, method, params, ctx)

	handlers := middleware.GetHandlersBeforeResponse()

	for _, handler := range handlers {
		handler(ctx)
	}
	return
}
