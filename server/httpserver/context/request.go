package context

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Request struct {
	ctx     *fiber.Ctx
	request *fasthttp.Request
}

func (this *Request) Method() string {
	return string(this.request.Header.Method())
}

func (this *Request) RequestURI() string {
	return string(this.request.Header.RequestURI())
}

func (this *Request) Header(key string, defaults ...string) string {
	return this.ctx.Get(key, defaults...)
}

func (this *Request) Headers() map[string]string {
	return this.ctx.GetReqHeaders()
}

func (this *Request) Path() string {
	return this.ctx.Path()
}

func (this *Request) Clear() {
	this.ctx = nil
	this.request = nil
}
