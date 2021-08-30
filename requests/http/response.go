package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Header http.Header
	Payload []byte
	StatusCode int
	StatusMessage string
	Raw *http.Response
}

func (this *Response) From (httpResponse *http.Response) (err error) {
	this.Header = httpResponse.Header
	this.Payload, err = ioutil.ReadAll(httpResponse.Body)

	if err != nil {
		return
	}

	this.StatusCode = httpResponse.StatusCode
	this.StatusMessage = httpResponse.Status
	this.Raw = httpResponse

	return
}

func (this *Response) BindJSON (v interface{}) (err error) {
	err = json.Unmarshal(this.Payload, v)
	
	return
}