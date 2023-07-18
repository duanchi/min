package route

import (
	"github.com/duanchi/min/server/handler"
	"github.com/duanchi/min/server/httpserver"
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/duanchi/min/server/middleware"
	"github.com/duanchi/min/server/types"
	"strings"
)

var RestfulRoutes = types.RestfulRoutesMap{}

func RestfulRoutesInit(httpServer *httpserver.Httpserver) {
	for key, _ := range RestfulRoutes {

		resource := key

		handlers := []httpserver.Handler{
			func(ctx *context.Context) error {
				ctx.Set("resource", resource)
				return ctx.Next()
			},
		}

		handlers = middleware.GetHandlersAfterRouteAppend(handlers)

		handlers = append(handlers, func(ctx *context.Context) error {
			handler.RestfulHandle(resource, RestfulRoutes[resource], ctx, httpServer)
			afterResponseHandlers := middleware.GetHandlersAfterResponse()
			if len(afterResponseHandlers) > 0 {
				go func() {
					for _, afterResponseHandler := range afterResponseHandlers {
						afterResponseHandler(ctx)
					}
				}()
			}
			return ctx.Next()
		})

		handlers = append(handlers, func(c *context.Context) error {
			c.Clear()
			return nil
		})

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
