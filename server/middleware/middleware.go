package middleware

import (
	_interface "github.com/duanchi/min/interface"
	"github.com/gin-gonic/gin"
	"reflect"
)

const (
	BeforeRoute = "beforeRoute"
	AfterRoute = "afterRoute"
	BeforeResponse = "beforeResponse"
	AfterResponse = "afterResponse"
	AfterPanic = "afterPanic"
)

var Middlewares []reflect.Value

/**
初始化before-route的中间件
 */
func Init (httpServer *gin.Engine, aop string) {
	for key, _ := range Middlewares {

		index := key

		switch aop {
		case BeforeRoute:
			httpServer.Use(Middlewares[index].Interface().(_interface.MiddlewareInterface).BeforeRoute)
		case AfterRoute:
			httpServer.Use(Middlewares[index].Interface().(_interface.MiddlewareInterface).AfterRoute)
		case BeforeResponse:
			httpServer.Use(Middlewares[index].Interface().(_interface.MiddlewareInterface).BeforeResponse)
		case AfterResponse:
			httpServer.Use(Middlewares[index].Interface().(_interface.MiddlewareInterface).AfterResponse)
		case AfterPanic:
			httpServer.Use(Middlewares[index].Interface().(_interface.MiddlewareInterface).AfterPanic)
		}

	}
}

func GetHandlersBeforeResponse () []gin.HandlerFunc {
	var handlers []gin.HandlerFunc
	for key, _ := range Middlewares {
		index := key
		handlers = append(handlers, Middlewares[index].Interface().(_interface.MiddlewareInterface).BeforeResponse)
	}

	return handlers
}

func GetHandlersAfterResponse () []gin.HandlerFunc {
	var handlers []gin.HandlerFunc
	for key, _ := range Middlewares {
		index := key
		handlers = append(handlers, Middlewares[index].Interface().(_interface.MiddlewareInterface).AfterResponse)
	}

	return handlers
}

func HandleAfterRoute (ctx *gin.Context) {
	for key, _ := range Middlewares {
		index := key
		Middlewares[index].Interface().(_interface.MiddlewareInterface).AfterRoute(ctx)
	}
}

func GetHandlersAfterRouter () []gin.HandlerFunc {
	var handlers []gin.HandlerFunc
	for key, _ := range Middlewares {
		index := key
		handlers = append(handlers, Middlewares[index].Interface().(_interface.MiddlewareInterface).AfterRoute)
	}

	return handlers
}

func  GetHandlersAfterRouterAppend (handlers []gin.HandlerFunc) []gin.HandlerFunc {
	for key, _ := range Middlewares {
		index := key
		handlers = append(handlers, Middlewares[index].Interface().(_interface.MiddlewareInterface).AfterRoute)
	}

	return handlers
}