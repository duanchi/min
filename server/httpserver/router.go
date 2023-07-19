package httpserver

import "github.com/gofiber/fiber/v2"

type Router interface {
	Use(args ...interface{}) Router

	GET(path string, handlers ...Handler) Router
	HEAD(path string, handlers ...Handler) Router
	POST(path string, handlers ...Handler) Router
	PUT(path string, handlers ...Handler) Router
	DELETE(path string, handlers ...Handler) Router
	CONNECT(path string, handlers ...Handler) Router
	OPTIONS(path string, handlers ...Handler) Router
	TRACE(path string, handlers ...Handler) Router
	PATCH(path string, handlers ...Handler) Router

	Add(method, path string, handlers ...Handler) Router
	Static(prefix, root string, config ...fiber.Static) Router
	Any(path string, handlers ...Handler) Router

	Group(prefix string, handlers ...Handler) Router

	Route(prefix string, fn func(router Router), name ...string) Router

	Mount(prefix string, fiber *Httpserver) Router

	Name(name string) Router
}
