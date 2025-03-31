package httpserver

import (
	"github.com/duanchi/min/v2/server/types"
	"github.com/gofiber/fiber/v2"
)

type Router interface {
	Use(args ...interface{}) Router

	GET(path string, handlers ...types.ServerHandleFunc) Router
	HEAD(path string, handlers ...types.ServerHandleFunc) Router
	POST(path string, handlers ...types.ServerHandleFunc) Router
	PUT(path string, handlers ...types.ServerHandleFunc) Router
	DELETE(path string, handlers ...types.ServerHandleFunc) Router
	CONNECT(path string, handlers ...types.ServerHandleFunc) Router
	OPTIONS(path string, handlers ...types.ServerHandleFunc) Router
	TRACE(path string, handlers ...types.ServerHandleFunc) Router
	PATCH(path string, handlers ...types.ServerHandleFunc) Router
	Add(method, path string, handlers ...types.ServerHandleFunc) Router
	Static(prefix, root string, config ...fiber.Static) Router
	ALL(path string, handlers ...types.ServerHandleFunc) Router
	Group(prefix string, handlers ...types.ServerHandleFunc) Router
	Route(prefix string, fn func(router Router), name ...string) Router
	Mount(prefix string, fiber *Httpserver) Router
	Name(name string) Router
}
