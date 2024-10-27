package middleware

import (
	"fmt"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/server/httpserver"
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/types/middleware"
	"github.com/duanchi/min/util"
	"reflect"
	"regexp"
	"strings"
)

const (
	BEFORE_ROUTE    = "beforeRoute"
	AFTER_ROUTE     = "afterRoute"
	BEFORE_RESPONSE = "beforeResponse"
	AFTER_RESPONSE  = "afterResponse"
	AFTER_PANIC     = "afterPanic"
)

var Middlewares []reflect.Value

var beforeRouteMiddlewares []interface{}
var afterRouteMiddlewares []types.ServerHandleFunc
var beforeResponseMiddlewares []types.ServerHandleFunc
var afterResponseMiddlewares []types.ServerHandleFunc
var afterPanicMiddlewares []types.ServerHandleFunc

/*
*
初始化before-route的中间件
*/
func Init(httpServer *httpserver.Httpserver) {
	for key, _ := range Middlewares {

		index := key
		middlewareBean := Middlewares[index].Interface().(_interface.MiddlewareInterface)

		beforeRouteMiddlewares = append(beforeRouteMiddlewares, types.ServerHandleFunc(middlewareBean.BeforeRoute))

		afterRouteMiddlewares = append(afterRouteMiddlewares, types.ServerHandleFunc(func(context *context.Context) (err error) {
			if matchRoute(middlewareBean.Includes(), middlewareBean.Excludes(), context) {
				return middlewareBean.AfterRoute(context)
			}
			if err != nil {
				return err
			} else {
				return context.Next()
			}
		}))

		beforeResponseMiddlewares = append(beforeResponseMiddlewares, types.ServerHandleFunc(func(context *context.Context) (err error) {
			if matchRoute(middlewareBean.Includes(), middlewareBean.Excludes(), context) {
				err = middlewareBean.BeforeResponse(context)
			}
			if err != nil {
				return err
			} else {
				return context.Next()
			}
		}))

		afterResponseMiddlewares = append(afterResponseMiddlewares, types.ServerHandleFunc(func(context *context.Context) (err error) {
			if matchRoute(middlewareBean.Includes(), middlewareBean.Excludes(), context) {
				err = middlewareBean.AfterResponse(context)
			}
			if err != nil {
				return err
			} else {
				return context.Next()
			}
		}))

		afterPanicMiddlewares = append(afterPanicMiddlewares, types.ServerHandleFunc(func(context *context.Context) (err error) {
			if matchRoute(middlewareBean.Includes(), middlewareBean.Excludes(), context) {
				err = middlewareBean.AfterPanic(context)
			}
			if err != nil {
				return err
			} else {
				return context.Next()
			}
		}))
	}

	fmt.Println("afterRouteMiddlewares", afterRouteMiddlewares)

	httpServer.Use(beforeRouteMiddlewares...)
}

func GetAfterRouteMiddlewares() []types.ServerHandleFunc {
	return afterRouteMiddlewares
}

func GetAfterResponseMiddlewares() []types.ServerHandleFunc {
	return afterResponseMiddlewares
}

func matchRoute(includes middleware.Includes, excludes middleware.Excludes, ctx *context.Context) bool {

	if includes != nil && len(includes) > 0 {
		if !match(includes, ctx) {
			return false
		}
	}

	if excludes != nil && len(excludes) > 0 {
		if match(excludes, ctx) {
			return false
		}
	}

	return true
}

func match(patterns map[string]string, ctx *context.Context) bool {
	for pattern, methods := range patterns {
		hasMethod := false
		methods = strings.ToUpper(methods)
		patternStack := strings.SplitN(pattern, ":", 2)
		if len(patternStack) == 1 {
			patternStack = append([]string{""}, patternStack[0])
		}
		if strings.Contains(methods, "ALL") {
			hasMethod = true
		} else if strings.Contains(methods, "WEBSOCKET") {
			/*upgradeRequest := ctx.Request().Header.Get("Connection")
			upgradeProtocol := ctx.Request.Header.Get("Upgrade")

			if upgradeRequest == "Upgrade" && strings.ToUpper(upgradeProtocol) == "WEBSOCKET" {
				hasMethod = true
			}*/
		} else {
			methodsStack := strings.Split(methods, ",")
			for _, method := range methodsStack {
				if s := strings.TrimSpace(method); s == string(ctx.Request().Method()) {
					hasMethod = true
					break
				}
			}
		}

		if hasMethod {
			switch patternStack[0] {
			case "":
				if strings.ContainsAny(patternStack[1], "*?[]!") {
					// fnmatch匹配
					if util.Fnmatch(patternStack[1], string(ctx.Request().RequestURI()), 0) {
						return true
					}
				} else {
					// 默认任意匹配
					if strings.Contains(string(ctx.Request().RequestURI()), patternStack[1]) {
						return true
					}
				}
			case "=":
				// 默认完全匹配
				if string(ctx.Request().RequestURI()) == patternStack[1] {
					return true
				}
			case "^":
				// 默认prefix匹配
				if strings.HasPrefix(string(ctx.Request().RequestURI()), patternStack[1]) {
					return true
				}
			case "~":
				regex := regexp.MustCompile(patternStack[1])
				if regex.MatchString(string(ctx.Request().RequestURI())) {
					return true
				}
			}
		}
	}
	return false
}

func GetHandlersBeforeResponse() []httpserver.Handler {
	var handlers []httpserver.Handler
	for key, _ := range Middlewares {
		index := key
		appendMiddleware := Middlewares[index].Interface().(_interface.MiddlewareInterface)
		handlers = append(handlers, func(context *context.Context) error {
			if matchRoute(appendMiddleware.Includes(), appendMiddleware.Excludes(), context) {
				appendMiddleware.BeforeResponse(context)
			}
			return nil
		})
	}

	return handlers
}

func GetHandlersAfterResponse() []httpserver.Handler {
	var handlers []httpserver.Handler
	for key, _ := range Middlewares {
		index := key
		appendMiddleware := Middlewares[index].Interface().(_interface.MiddlewareInterface)
		handlers = append(handlers, func(context *context.Context) error {
			if matchRoute(appendMiddleware.Includes(), appendMiddleware.Excludes(), context) {
				appendMiddleware.AfterResponse(context)
			}
			return nil
		})
	}

	return handlers
}

func HandleAfterRoute(ctx *context.Context) error {
	for key, _ := range Middlewares {
		index := key
		appendMiddleware := Middlewares[index].Interface().(_interface.MiddlewareInterface)
		func(context *context.Context) {
			if matchRoute(appendMiddleware.Includes(), appendMiddleware.Excludes(), context) {
				appendMiddleware.AfterRoute(context)
			}
		}(ctx)
	}
	return nil
}

func GetHandlersAfterRoute() []httpserver.Handler {
	var handlers []httpserver.Handler
	for key, _ := range Middlewares {
		index := key
		appendMiddleware := Middlewares[index].Interface().(_interface.MiddlewareInterface)
		handlers = append(handlers, func(context *context.Context) error {
			if matchRoute(appendMiddleware.Includes(), appendMiddleware.Excludes(), context) {
				appendMiddleware.AfterRoute(context)
			}
			return nil
		})
	}

	return handlers
}

func GetHandlersAfterRouteAppend(handlers []httpserver.Handler) []httpserver.Handler {
	for key, _ := range Middlewares {
		index := key
		appendMiddleware := Middlewares[index].Interface().(_interface.MiddlewareInterface)
		handlers = append(handlers, func(context *context.Context) error {
			if matchRoute(appendMiddleware.Includes(), appendMiddleware.Excludes(), context) {
				appendMiddleware.AfterRoute(context)
			}
			return nil
		})
	}

	return handlers
}
