package route

import (
	"github.com/duanchi/min/server/handler"
	"github.com/duanchi/min/server/httpserver"
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/duanchi/min/server/middleware"
	"github.com/duanchi/min/server/types"

	"strings"
)

var BaseRoutes = types.BaseRoutesMap{}

func BaseRouteInit(httpServer *httpserver.Httpserver) {
	afterResponseMiddlewares := middleware.GetAfterResponseMiddlewares()
	afterRouteMiddlewares := middleware.GetAfterRouteMiddlewares()

	for key, route := range BaseRoutes {
		methods := []string{"ALL"}
		if route.Method != "" {
			methods = strings.Split(strings.ToUpper(route.Method), ",")
		}

		handleBeanKey := key
		handlers := []types.ServerHandleFunc{}
		if len(afterRouteMiddlewares) > 0 {
			handlers = append(handlers, afterRouteMiddlewares...)
		}
		handlers = append(
			handlers,
			func(ctx *context.Context) {
				handler.RouteHandle(route.Path, BaseRoutes[handleBeanKey].Value, ctx, httpServer)
			},
		)

		if len(afterResponseMiddlewares) > 0 {
			handlers = append(handlers, afterResponseMiddlewares...)
		}

		for _, method := range methods {
			httpServer.Add(method, route.Path, handlers...)
		}
	}
}
