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

		// name := key

		// stack := strings.SplitN(name, "@", 2)
		// route := "/"
		methods := []string{"ALL"}

		//if stack[0] != "" {
		//	route = stack[0]
		//}
		if route.Method != "" {
			methods = strings.Split(strings.ToUpper(route.Method), ",")
		}

		handleBeanKey := key

		for _, method := range methods {

			handlers := []types.ServerHandleFunc{
				func(ctx *context.Context) {
					if len(afterRouteMiddlewares) > 0 {
						for _, handler := range afterRouteMiddlewares {
							if ctx.IsNext() {
								handler(ctx)
							}
							return
						}
					}
				},
				func(ctx *context.Context) {
					handler.RouteHandle(route.Path, BaseRoutes[handleBeanKey].Value, ctx, httpServer)
					if len(afterResponseMiddlewares) > 0 {
						go func() {
							for _, afterResponseMiddleware := range afterResponseMiddlewares {
								afterResponseMiddleware(ctx)
							}
						}()
					}
					ctx.Clear()
				},
			}

			httpServer.Add(method, route.Path, handlers...)
		}
	}
}
