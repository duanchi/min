package route

import (
	"github.com/duanchi/min/server/handler"
	"github.com/duanchi/min/server/httpserver"
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/duanchi/min/server/middleware"
	server_types "github.com/duanchi/min/server/types"
	"github.com/duanchi/min/types"
	"strings"
)

var RestfulRoutes = server_types.RestfulRoutesMap{}

func RestfulRoutesInit(httpServer *httpserver.Httpserver) {
	for key, _ := range RestfulRoutes {

		resource := key

		handlers := append([]types.ServerHandleFunc{
			func(ctx *context.Context) error {
				ctx.Set("resource", resource)
				return ctx.Next()
			},
		}, middleware.GetAfterRouteMiddlewares()...)

		// handlers = middleware.GetHandlersAfterRouteAppend(handlers)

		handlers = append(handlers,
			func(ctx *context.Context) error {
				return handler.RestfulHandle(resource, RestfulRoutes[resource], ctx, httpServer)
			},
			func(ctx *context.Context) error {
				afterResponseMiddlewares := middleware.GetAfterResponseMiddlewares()
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

		if strings.Contains(resource, ":id") {
			resource := strings.ReplaceAll("/"+resource, "//", "/")
			httpServer.Any(resource, handlers...)
			if !strings.HasSuffix(resource, "/") {
				httpServer.Any(resource+"/", handlers...)
			}
		} else {
			httpServer.Any("/"+resource, handlers...)
			httpServer.Any("/"+resource+"/", handlers...)
			httpServer.Any("/"+resource+"/:id", handlers...)
		}
	}
}
