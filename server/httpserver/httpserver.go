package httpserver

import (
	"fmt"
	"github.com/duanchi/min/v2/server/httpserver/context"
	"github.com/duanchi/min/v2/server/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"reflect"
)

type Httpserver struct {
	instance *fiber.App
}

func New(config interface{}) *Httpserver {
	app := fiber.New()
	app.Use(logger.New(logger.Config{Format: "[${ip}]:${port} ${status} - ${method} ${path}\n"}))
	return &Httpserver{
		instance: app,
	}
}

func (this *Httpserver) Instance() *fiber.App {
	return this.instance
}

func (this *Httpserver) Listen(host string, port string) error {
	err := this.instance.Listen(host + ":" + port)
	return err
}

func (this *Httpserver) SetLogLevel(level int) {
	switch level {
	case LOG_TRACE:
		log.SetLevel(log.LevelTrace)
	case LOG_DEBUG:
		log.SetLevel(log.LevelDebug)
	case LOG_INFO:
		log.SetLevel(log.LevelInfo)
	case LOG_WARN:
		log.SetLevel(log.LevelWarn)
	case LOG_ERROR:
		log.SetLevel(log.LevelError)
	case LOG_FATAL:
		log.SetLevel(log.LevelFatal)
	case LOG_PANIC:
		log.SetLevel(log.LevelPanic)
	}
}

// Add allows you to specify a HTTP method to register a route
func (this *Httpserver) Add(method, path string, handlers ...types.ServerHandleFunc) Router {
	fiberHandlers := []fiber.Handler{}
	for _, h := range handlers {
		fiberHandlers = append(fiberHandlers, func(ctx *fiber.Ctx) error {
			h(ctx.Locals("CONTEXT", context.New(ctx)).(*context.Context))
			return nil
		})
	}
	if method == METHOD_ALL {
		this.instance.All(path, fiberHandlers...)
	} else {
		this.instance.Add(method, path, fiberHandlers...)
	}
	return this
}

// Static will create a file server serving static files
func (this *Httpserver) Static(prefix, root string, config ...fiber.Static) Router {
	this.instance.Static(prefix, root, config...)
	return this
}

func (this *Httpserver) Use(args ...interface{}) Router {
	var prefix string
	var prefixes []string
	var handlers []types.ServerHandleFunc

	for i := 0; i < len(args); i++ {
		switch arg := args[i].(type) {
		case string:
			prefix = arg
		case []string:
			prefixes = arg
		case types.ServerHandleFunc:
			handlers = append(handlers, arg)
		default:
			panic(fmt.Sprintf("use: invalid handler %v\n", reflect.TypeOf(arg)))
		}
	}

	if len(prefixes) == 0 {
		prefixes = append(prefixes, prefix)
	}

	for _, prefix := range prefixes {
		this.Add(METHOD_USE, prefix, handlers...)
	}

	return this
}

func (this *Httpserver) GET(path string, handlers ...types.ServerHandleFunc) Router {
	return this.HEAD(path, handlers...).Add(METHOD_GET, path, handlers...)
}
func (this *Httpserver) HEAD(path string, handlers ...types.ServerHandleFunc) Router {
	return this.Add(METHOD_HEAD, path, handlers...)
}
func (this *Httpserver) POST(path string, handlers ...types.ServerHandleFunc) Router {
	return this.Add(METHOD_POST, path, handlers...)
}
func (this *Httpserver) PUT(path string, handlers ...types.ServerHandleFunc) Router {
	return this.Add(METHOD_PUT, path, handlers...)
}
func (this *Httpserver) DELETE(path string, handlers ...types.ServerHandleFunc) Router {
	return this.Add(METHOD_DELETE, path, handlers...)
}
func (this *Httpserver) CONNECT(path string, handlers ...types.ServerHandleFunc) Router {
	return this.Add(METHOD_CONNECT, path, handlers...)
}
func (this *Httpserver) OPTIONS(path string, handlers ...types.ServerHandleFunc) Router {
	return this.Add(METHOD_OPTIONS, path, handlers...)
}
func (this *Httpserver) TRACE(path string, handlers ...types.ServerHandleFunc) Router {
	return this.Add(METHOD_TRACE, path, handlers...)
}
func (this *Httpserver) PATCH(path string, handlers ...types.ServerHandleFunc) Router {
	return this.Add(METHOD_PATCH, path, handlers...)
}

func (this *Httpserver) ALL(path string, handlers ...types.ServerHandleFunc) Router {
	return this.Add(METHOD_ALL, path, handlers...)
}

func (this *Httpserver) Group(prefix string, handlers ...types.ServerHandleFunc) Router {
	return this
}

func (this *Httpserver) Route(prefix string, fn func(router Router), name ...string) Router {
	return this
}

func (this *Httpserver) Mount(prefix string, fiber *Httpserver) Router {
	return this
}

func (this *Httpserver) Name(name string) Router {
	return this
}

func (this *Httpserver) Stop() error {
	return this.instance.Shutdown()
}

func NewContext(ctx *fiber.Ctx) *context.Context {
	return context.New(ctx)
}
