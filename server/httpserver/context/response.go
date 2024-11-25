package context

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
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
		this.response.Header.SetStatusMessage([]byte(utils.StatusMessage(code)))
	} else {
		this.response.Header.SetStatusMessage([]byte(message[0]))
	}

	return this
}

func (this *Response) Clear() {
	this.ctx = nil
	this.response = nil
}
