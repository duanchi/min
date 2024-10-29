package route

import (
	"github.com/duanchi/min/server/handler"
	"github.com/duanchi/min/server/httpserver"
	"github.com/duanchi/min/server/httpserver/constant"
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

		handlers := []types.ServerHandleFunc{
			func(ctx *context.Context) {
				ctx.Set(constant.RESOURCE, resource)
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
				handler.RestfulHandle(resource, RestfulRoutes[resource], ctx, httpServer)
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
