package context

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"strings"
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

func (this *Request) Headers() Header {
	return this.ctx.GetReqHeaders()
}

func (this *Request) Path() string {
	return this.ctx.Path()
}

func (this *Request) Query(key string, defaults ...string) string {
	return this.ctx.Query(key, defaults...)
}

func (this *Request) Body() []byte {
	return this.ctx.Body()
}

func (this *Request) RemoteAddr() string {
	IPAddress := this.ctx.GetRespHeader("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = this.ctx.GetRespHeader("X-Forwarded-For")
	}
	return strings.TrimSpace(strings.Split(IPAddress, ",")[0])
}

func (this *Request) Clear() {
	this.ctx = nil
	this.request = nil
}
