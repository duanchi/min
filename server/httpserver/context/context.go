package context

import (
	"github.com/duanchi/min/server/httpserver/constant"
	"github.com/duanchi/min/server/validate"
	"github.com/gofiber/fiber/v2"
)

type Context struct {
	ctx      *fiber.Ctx
	request  *Request
	response *Response
	params   *Params
	next     bool
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
		next: true,
	}
}

func (this *Context) Ctx() *fiber.Ctx {
	return this.ctx
}

func (this *Context) Request() *Request {
	return this.request
}

func (this *Context) Response() *Response { return this.response }

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

func (this *Context) Has(key string) bool {
	return this.ctx.Locals(key) != nil
}

func (this *Context) Query(key string, defaults ...string) string {
	return this.request.ctx.Query(key, defaults...)
}

func (this *Context) GetHeader(key string) string {
	return this.ctx.Get(key)
}

func (this *Context) Next() {
	this.next = true
	return
}

func (this *Context) JSON(obj interface{}) error {
	this.ctx.Response().SetStatusCode(constant.StatusOK)
	this.next = false
	return this.ctx.JSON(obj)
}

func (this *Context) JSONWithStatus(code int, obj interface{}) error {
	this.ctx.Response().SetStatusCode(code)
	this.next = false
	return this.ctx.JSON(obj)
}

func (this *Context) Abort() {
	this.next = false
}

func (this *Context) AbortWithStatus(code int) error {
	this.next = false
	return this.ctx.Status(code).SendString("")
}

func (this *Context) DataWithStatus(code int, data []byte) error {
	return this.ctx.Status(code).Send(data)
}

func (this *Context) IsNext() bool {
	return this.next
}

func (this *Context) Bind(obj interface{}) error {
	result := this.ctx.BodyParser(&obj)

	err := validate.Validate(obj)
	if err != nil {
		return err
	}
	return result
}

func (this *Context) Clear() {
	this.request.Clear()
	this.response.Clear()
	this.ctx = nil
	this.request = nil
	this.response = nil
}
