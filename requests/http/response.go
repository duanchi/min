package http

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type Response struct {
	Header        fasthttp.ResponseHeader
	Payload       []byte
	StatusCode    int
	StatusMessage string
	Raw           *fasthttp.Response
}

func (this *Response) From(httpResponse *fasthttp.Response) (err error) {
	httpResponse.Header.CopyTo(&this.Header)
	this.Payload = httpResponse.Body()
	this.StatusCode = httpResponse.StatusCode()
	this.StatusMessage = string(httpResponse.Header.StatusMessage())
	this.Raw = httpResponse

	return
}

func (this *Response) BindJSON(v interface{}) (err error) {
	err = json.Unmarshal(this.Payload, v)

	return
}
