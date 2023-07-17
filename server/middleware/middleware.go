package middleware

import (
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/types/middleware"
	"github.com/duanchi/min/util"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"regexp"
	"strings"
)

const (
	BeforeRoute    = "beforeRoute"
	AfterRoute     = "afterRoute"
	BeforeResponse = "beforeResponse"
	AfterResponse  = "afterResponse"
	AfterPanic     = "afterPanic"
)

var Middlewares []reflect.Value

/*
*
初始化before-route的中间件
*/
func Init(httpServer *fiber.App, aop string) {
	for key, _ := range Middlewares {

		index := key
		middleware := Middlewares[index].Interface().(_interface.MiddlewareInterface)
		switch aop {
		case BeforeRoute:
			httpServer.Use(middleware.BeforeRoute)
		case AfterRoute:
			httpServer.Use(func(context *fiber.Ctx) {
				if matchRoute(middleware.Includes(), middleware.Excludes(), context) {
					middleware.AfterRoute(context)
				}
			})
		case BeforeResponse:
			httpServer.Use(func(context *fiber.Ctx) {
				if matchRoute(middleware.Includes(), middleware.Excludes(), context) {
					middleware.BeforeResponse(context)
				}
			})
		case AfterResponse:
			httpServer.Use(func(context *fiber.Ctx) {
				if matchRoute(middleware.Includes(), middleware.Excludes(), context) {
					middleware.AfterResponse(context)
				}
			})
		case AfterPanic:
			httpServer.Use(func(context *fiber.Ctx) {
				if matchRoute(middleware.Includes(), middleware.Excludes(), context) {
					middleware.AfterPanic(context)
				}
			})
		}

	}
}

func matchRoute(includes middleware.Includes, excludes middleware.Excludes, ctx *fiber.Ctx) bool {

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

func match(patterns map[string]string, ctx *fiber.Ctx) bool {
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
				if s := strings.TrimSpace(method); s == string(ctx.Request().Header.Method()) {
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
					if util.Fnmatch(patternStack[1], string(ctx.Request().Header.RequestURI()), 0) {
						return true
					}
				} else {
					// 默认任意匹配
					if strings.Contains(string(ctx.Request().Header.RequestURI()), patternStack[1]) {
						return true
					}
				}
			case "=":
				// 默认完全匹配
				if string(ctx.Request().Header.RequestURI()) == patternStack[1] {
					return true
				}
			case "^":
				// 默认prefix匹配
				if strings.HasPrefix(string(ctx.Request().Header.RequestURI()), patternStack[1]) {
					return true
				}
			case "~":
				regex := regexp.MustCompile(patternStack[1])
				if regex.MatchString(string(ctx.Request().Header.RequestURI())) {
					return true
				}
			}
		}
	}
	return false
}

func GetHandlersBeforeResponse() []fiber.Handler {
	var handlers []fiber.Handler
	for key, _ := range Middlewares {
		index := key
		appendMiddleware := Middlewares[index].Interface().(_interface.MiddlewareInterface)
		handlers = append(handlers, func(context *fiber.Ctx) error {
			if matchRoute(appendMiddleware.Includes(), appendMiddleware.Excludes(), context) {
				appendMiddleware.BeforeResponse(context)
			}
			return nil
		})
	}

	return handlers
}

func GetHandlersAfterResponse() []fiber.Handler {
	var handlers []fiber.Handler
	for key, _ := range Middlewares {
		index := key
		appendMiddleware := Middlewares[index].Interface().(_interface.MiddlewareInterface)
		handlers = append(handlers, func(context *fiber.Ctx) error {
			if matchRoute(appendMiddleware.Includes(), appendMiddleware.Excludes(), context) {
				appendMiddleware.AfterResponse(context)
			}
			return nil
		})
	}

	return handlers
}

func HandleAfterRoute(ctx *fiber.Ctx) {
	for key, _ := range Middlewares {
		index := key
		appendMiddleware := Middlewares[index].Interface().(_interface.MiddlewareInterface)
		func(context *fiber.Ctx) {
			if matchRoute(appendMiddleware.Includes(), appendMiddleware.Excludes(), context) {
				appendMiddleware.AfterRoute(context)
			}
		}(ctx)
	}
}

func GetHandlersAfterRouter() []fiber.Handler {
	var handlers []fiber.Handler
	for key, _ := range Middlewares {
		index := key
		appendMiddleware := Middlewares[index].Interface().(_interface.MiddlewareInterface)
		handlers = append(handlers, func(context *fiber.Ctx) error {
			if matchRoute(appendMiddleware.Includes(), appendMiddleware.Excludes(), context) {
				appendMiddleware.AfterRoute(context)
			}
			return nil
		})
	}

	return handlers
}

func GetHandlersAfterRouterAppend(handlers []fiber.Handler) []fiber.Handler {
	for key, _ := range Middlewares {
		index := key
		appendMiddleware := Middlewares[index].Interface().(_interface.MiddlewareInterface)
		handlers = append(handlers, func(context *fiber.Ctx) error {
			if matchRoute(appendMiddleware.Includes(), appendMiddleware.Excludes(), context) {
				appendMiddleware.AfterRoute(context)
			}
			return nil
		})
	}

	return handlers
}
