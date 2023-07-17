package handler

import (
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/server/middleware"
	"github.com/gofiber/fiber/v2"
	"reflect"
)

func RouteHandle(path string, handle reflect.Value, ctx *fiber.Ctx, engine *fiber.App) {
	params := ctx.Params
	method := string(ctx.Request().Header.Method())

	handle.Interface().(_interface.RouterInterface).Handle(string(ctx.Request().URI().Path()), method, params, ctx)

	handlers := middleware.GetHandlersBeforeResponse()

	for _, handler := range handlers {
		handler(ctx)
	}
	return
}
