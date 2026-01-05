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

func (this *Response) GetHeader(key string) string {
	return string(this.Header.Peek(key))
}

func (this *Response) GetHeaders() map[string]string {
	keys := this.Header.PeekKeys()
	maps := make(map[string]string, len(keys))
	for _, keyBytes := range keys {
		key := string(keyBytes)
		maps[key] = string(this.Header.Peek(key))
	}
	return maps
}

func (this *Response) BindJSON(v interface{}) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	err = json.Unmarshal(this.Payload, v)

	return
}
