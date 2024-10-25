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

		for _, method := range methods {

			handlers := append(
				afterRouteMiddlewares,
				func(ctx *context.Context) error {
					return handler.RouteHandle(route.Path, BaseRoutes[key].Value, ctx, httpServer)
				},
				func(ctx *context.Context) error {
					if len(afterResponseMiddlewares) > 0 {
						go func() {
							for _, afterResponseMiddleware := range afterResponseMiddlewares {
								afterResponseMiddleware(ctx)
							}
						}()
					}
					return ctx.Next()
				},
				func(c *context.Context) error {
					c.Clear()
					return nil
				},
			)

			if method == "ALL" {
				httpServer.ALL(route.Path, handlers...)
			} else {
				httpServer.Add(method, route.Path, handlers...)
			}
		}
	}
}
