package route

import (
	"github.com/duanchi/min/server/handler"
	"github.com/duanchi/min/server/middleware"
	"github.com/gin-gonic/gin"
	"reflect"
)

type RestfulRoutesMap map[string]reflect.Value

var RestfulRoutes = RestfulRoutesMap{}

func (this RestfulRoutesMap) Init (httpServer *gin.Engine) {
	for key, _ := range this {

		resource := key

		handlers := []gin.HandlerFunc{
			func (ctx *gin.Context) {
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
					for _, afterResponseHandler  := range afterResponseHandlers {
						afterResponseHandler(ctx)
					}
				}()
			}
			ctx.Next()
		})

		httpServer.Any("/" + resource, handlers...)
		httpServer.Any("/" + resource + "/", handlers...)
		httpServer.Any("/" + resource + "/:id", handlers...)
	}
}