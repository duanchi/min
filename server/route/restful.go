package route

import (
	"github.com/duanchi/min/server/handler"
	"github.com/duanchi/min/server/httpserver"
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/duanchi/min/server/middleware"
	serverTypes "github.com/duanchi/min/server/types"
	"github.com/duanchi/min/types"
	"strings"
)

var RestfulRoutes = serverTypes.RestfulRoutesMap{}

func RestfulRouteInit(httpServer *httpserver.Httpserver) {
	afterResponseMiddlewares := middleware.GetAfterResponseMiddlewares()
	afterRouteMiddlewares := middleware.GetAfterRouteMiddlewares()

	for key, route := range RestfulRoutes {

		resource := key

		handlers := append([]types.ServerHandleFunc{
			func(ctx *context.Context) error {
				ctx.Set("resource", resource)
				return ctx.Next()
			},
		}, afterRouteMiddlewares...)

		// handlers = middleware.GetHandlersAfterRouteAppend(handlers)

		handlers = append(handlers,
			func(ctx *context.Context) error {
				return handler.RestfulHandle(resource, RestfulRoutes[resource], ctx, httpServer)
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

		/*handlers = append(handlers, func(c *context.Context) error {
			c.Clear()
			return nil
		})*/

		if strings.Contains(resource, ":"+route.ResourceKey) {
			resource := strings.ReplaceAll("/"+resource, "//", "/")
			httpServer.ALL(resource, handlers...)
			if !strings.HasSuffix(resource, "/") {
				httpServer.ALL(resource+"/", handlers...)
			}
		} else {
			httpServer.ALL("/"+resource, handlers...)
			httpServer.ALL("/"+resource+"/", handlers...)
			httpServer.ALL("/"+resource+"/:"+route.ResourceKey, handlers...)
		}
	}
}
