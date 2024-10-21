package context

import (
	"github.com/duanchi/min/server/httpserver/constant"
	"github.com/gofiber/fiber/v2"
)

type Context struct {
	ctx      *fiber.Ctx
	request  *Request
	response *Response
	params   *Params
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
		params: &Params{
			params: ctx.AllParams(),
		},
	}
}

func (this *Context) Request() *Request {
	return this.request
}

func (this *Context) Param(key string, defaults ...string) string {
	return this.params.Get(key, defaults...)
}

func (this *Context) Params() *Params {
	return this.params
}

func (this *Context) Get(key string) interface{} {
	return this.ctx.Locals(key)
}

func (this *Context) Set(key string, value interface{}) {
	this.ctx.Locals(key, value)
}

func (this *Context) GetHeader(key string) string {
	return this.ctx.GetRespHeader(key)
}

func (this *Context) Next() error {
	return this.ctx.Next()
}

func (this *Context) JSON(obj interface{}) error {
	this.ctx.Response().SetStatusCode(constant.StatusOK)
	return this.ctx.JSON(obj)
}

func (this *Context) JSONWithStatus(code int, obj interface{}) error {
	this.ctx.Response().SetStatusCode(code)
	return this.ctx.JSON(obj)
}

func (this *Context) AbortWithStatus(code int) error {
	return this.ctx.Status(code).SendString("")
}

func (this *Context) Bind(obj interface{}) error {
	return this.ctx.BodyParser(&obj)
}

func (this *Context) Clear() {
	this.request.Clear()
	this.response.Clear()
	this.ctx = nil
	this.request = nil
	this.response = nil
}
