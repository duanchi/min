package context

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Response struct {
	ctx      *fiber.Ctx
	response *fasthttp.Response
}

func (this *Response) SetHeader(key string, value string) *Response {
	this.response.Header.Set(key, value)
	return this
}

func (this *Response) SetStatus(code int, message ...string) *Response {
	this.response.Header.SetStatusCode(code)
	if len(message) == 0 {
		this.response.Header.SetStatusMessage(fiber.StatusTeapot)
	}

}

func (this *Response) Clear() {
	this.ctx = nil
	this.response = nil
}
