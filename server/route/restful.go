package route

import (
	"github.com/duanchi/min/v2/server/handler"
	"github.com/duanchi/min/v2/server/httpserver"
	"github.com/duanchi/min/v2/server/httpserver/constant"
	"github.com/duanchi/min/v2/server/httpserver/context"
	"github.com/duanchi/min/v2/server/middleware"
	"github.com/duanchi/min/v2/server/types"
	"strings"
)

var RestfulRoutes = types.RestfulRoutesMap{}

func RestfulRouteInit(httpServer *httpserver.Httpserver) {
	afterResponseMiddlewares := middleware.GetAfterResponseMiddlewares()
	afterRouteMiddlewares := middleware.GetAfterRouteMiddlewares()

	for key, route := range RestfulRoutes {
		resource := strings.ReplaceAll("/"+key, "//", "/")

		handlers := []types.ServerHandleFunc{}
		if len(afterRouteMiddlewares) > 0 {
			handlers = append(handlers, afterRouteMiddlewares...)
		}
		handlers = append(
			handlers,
			func(ctx *context.Context) {
				ctx.Set(constant.RESOURCE, resource)
				ctx.Next()
			},
			func(ctx *context.Context) {
				handler.RestfulHandle(resource, RestfulRoutes[resource], ctx, httpServer)
			},
		)

		if len(afterResponseMiddlewares) > 0 {
			handlers = append(handlers, afterResponseMiddlewares...)
		}

		if strings.Contains(resource, ":"+route.ResourceKey) {
			httpServer.ALL(resource, handlers...)
			if !strings.HasSuffix(resource, "/") {
				httpServer.ALL(resource+"/", handlers...)
			}
		} else {
			httpServer.ALL("/"+resource, handlers...)
			// httpServer.ALL("/"+resource+"/", handlers...)
			httpServer.ALL("/"+resource+"/:"+route.ResourceKey+"?", handlers...)
		}
	}
}
