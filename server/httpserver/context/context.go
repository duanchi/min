package context

import (
	"github.com/duanchi/min/server/httpserver"
	"github.com/gofiber/fiber/v2"
)

type Context struct {
	ctx      *fiber.Ctx
	request  *Request
	response *Response
}

func New(ctx *fiber.Ctx) *Context {
	return &Context{
		ctx: ctx,
		request: &Request{
			ctx:     ctx,
			request: ctx.Request(),
		},
		response: &Response{
			ctx:      ctx,
			response: ctx.Response(),
		},
	}
}

func (this *Context) Request() *Request {
	return this.request
}

func (this *Context) Param(key string, defaults ...string) string {
	return this.ctx.Params(key, defaults...)
}

func (this *Context) Params() map[string]string {
	return this.ctx.AllParams()
}

func (this *Context) Get(key string) interface{} {
	return this.ctx.Locals(key)
}

func (this *Context) Set(key string, value interface{}) {
	this.ctx.Locals(key, value)
}

func (this *Context) Next() error {
	return this.ctx.Next()
}

func (this *Context) JSON(obj interface{}) error {
	this.ctx.Response().SetStatusCode(httpserver.StatusOK)
	return this.ctx.JSON(obj)
}

func (this *Context) JSONWithStatusCode(code int, obj interface{}) error {
	this.ctx.Response().SetStatusCode(code)
	return this.ctx.JSON(obj)
}

func (this *Context) Clear() {
	this.request.Clear()
	this.response.Clear()
	this.ctx = nil
	this.request = nil
	this.response = nil
}
