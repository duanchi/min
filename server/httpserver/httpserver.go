package httpserver

import (
	"fmt"
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"reflect"
)

type Httpserver struct {
	instance *fiber.App
}

func New(config interface{}) *Httpserver {
	return &Httpserver{
		instance: fiber.New(config.(fiber.Config)),
	}
}

func (this *Httpserver) Listen(host string, port string) error {
	return this.instance.Listen(host + ":" + port)
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
func (this *Httpserver) Add(method, path string, handlers ...Handler) Router {
	this.instance.Add(method, path, toFiberHandlers(handlers...)...)
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
	var handlers []Handler

	for i := 0; i < len(args); i++ {
		switch arg := args[i].(type) {
		case string:
			prefix = arg
		case []string:
			prefixes = arg
		case Handler:
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

func (this *Httpserver) Get(path string, handlers ...Handler) Router {
	return this.Head(path, handlers...).Add(METHOD_GET, path, handlers...)
}
func (this *Httpserver) Head(path string, handlers ...Handler) Router {
	return this.Add(METHOD_HEAD, path, handlers...)
}
func (this *Httpserver) Post(path string, handlers ...Handler) Router {
	return this.Add(METHOD_POST, path, handlers...)
}
func (this *Httpserver) Put(path string, handlers ...Handler) Router {
	return this.Add(METHOD_PUT, path, handlers...)
}
func (this *Httpserver) Delete(path string, handlers ...Handler) Router {
	return this.Add(METHOD_DELETE, path, handlers...)
}
func (this *Httpserver) Connect(path string, handlers ...Handler) Router {
	return this.Add(METHOD_CONNECT, path, handlers...)
}
func (this *Httpserver) Options(path string, handlers ...Handler) Router {
	return this.Add(METHOD_OPTIONS, path, handlers...)
}
func (this *Httpserver) Trace(path string, handlers ...Handler) Router {
	return this.Add(METHOD_TRACE, path, handlers...)
}
func (this *Httpserver) Patch(path string, handlers ...Handler) Router {
	return this.Add(METHOD_PATCH, path, handlers...)
}

func (this *Httpserver) Any(path string, handlers ...Handler) Router {
	for _, method := range []string{METHOD_GET, METHOD_POST, METHOD_PUT, METHOD_CONNECT, METHOD_DELETE, METHOD_OPTIONS, METHOD_HEAD, METHOD_PATCH, METHOD_TRACE} {
		_ = this.Add(method, path, handlers...)
	}
	return this
}

func (this *Httpserver) Group(prefix string, handlers ...Handler) Router {
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

func toFiberHandlers(handlers ...Handler) []fiber.Handler {
	fiberHandlers := []fiber.Handler{}
	for n, _ := range handlers {
		fiberHandlers = append(fiberHandlers, func(ctx *fiber.Ctx) error {
			return handlers[n](NewContext(ctx))
		})
	}
	return fiberHandlers
}
