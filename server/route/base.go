package route

import (
	"github.com/duanchi/min/server/handler"
	"github.com/duanchi/min/server/middleware"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
)

type BaseRoutesMap map[string]reflect.Value

var BaseRoutes = BaseRoutesMap{}

func (this BaseRoutesMap) Init(httpServer *fiber.App) {
	for key, _ := range this {

		name := key

		stack := strings.SplitN(name, "@", 2)
		route := "/"
		methods := []string{"ALL"}

		if stack[0] != "" {
			route = stack[0]
		}
		if len(stack) > 1 && stack[1] != "" {
			methods = strings.Split(strings.ToUpper(stack[1]), ",")
		}

		for _, method := range methods {

			handlers := middleware.GetHandlersAfterRouter()

			handlers = append(handlers, func(ctx *fiber.Handler) {
				handler.RouteHandle(route, BaseRoutes[name], ctx, httpServer)
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

			if method == "ALL" {
				httpServer.Any(route, handlers...)
			} else {
				httpServer.Handle(method, route, handlers...)
			}
		}
	}
}
