package route

import (
	"github.com/duanchi/min/server/handler"
	"github.com/duanchi/min/server/middleware"
	"github.com/duanchi/min/server/types"
	"github.com/gin-gonic/gin"
	"strings"
)

var RestfulRoutes = types.RestfulRoutesMap{}

func RestfulRoutesInit(httpServer *gin.Engine) {
	for key, _ := range RestfulRoutes {

		resource := key

		handlers := []gin.HandlerFunc{
			func(ctx *gin.Context) {
				ctx.Set("resource", resource)
				ctx.Next()
			},
		}

		handlers = middleware.GetHandlersAfterRouterAppend(handlers)

		handlers = append(handlers, func(ctx *gin.Context) {
			handler.RestfulHandle(resource, RestfulRoutes[resource], ctx, httpServer)
			afterResponseHandlers := middleware.GetHandlersAfterResponse()
			if len(afterResponseHandlers) > 0 {
				go func() {
					for _, afterResponseHandler := range afterResponseHandlers {
						afterResponseHandler(ctx)
					}
				}()
			}
			ctx.Next()
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
