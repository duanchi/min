package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/duanchi/min/util/arrays"
	"github.com/fatih/structs"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Request struct {
	initialed bool
	error     error

	method      string
	url         string
	baseUrl     string
	queryString string
	header      fasthttp.RequestHeader
	payload     []byte
	formData    interface{}
}

const (
	METHOD_POST    = "POST"
	METHOD_GET     = "GET"
	METHOD_PUT     = "PUT"
	METHOD_DELETE  = "DELETE"
	METHOD_HEAD    = "HEAD"
	METHOD_PATCH   = "PATCH"
	METHOD_OPTIONS = "OPTIONS"
	METHOD_TRACE   = "TRACE"
)

func POST(url string) *Request {
	request := New()

	request.Method(METHOD_POST)
	request.Url(url)

	return &request
}

func GET(url string) *Request {
	request := New()

	request.Method(METHOD_GET)
	request.Url(url)

	return &request
}

func PUT(url string) *Request {
	request := New()

	request.Method(METHOD_PUT)
	request.Url(url)

	return &request
}

func DELETE(url string) *Request {
	request := New()

	request.Method(METHOD_DELETE)
	request.Url(url)

	return &request
}

func HEAD(url string) *Request {
	request := New()

	request.Method(METHOD_HEAD)
	request.Url(url)

	return &request
}

func PATCH(url string) *Request {
	request := New()

	request.Method(METHOD_PATCH)
	request.Url(url)

	return &request
}

func OPTIONS(url string) *Request {
	request := New()

	request.Method(METHOD_OPTIONS)
	request.Url(url)

	return &request
}

func TRACE(url string) *Request {
	request := New()

	request.Method(METHOD_TRACE)
	request.Url(url)

	return &request
}

func (this *Request) Url(url string) *Request {
	this.url = url

	return this
}

func (this *Request) Method(method string) *Request {
	if _, has := arrays.ContainsString([]string{METHOD_GET, METHOD_POST, METHOD_PUT, METHOD_DELETE, METHOD_OPTIONS, METHOD_PATCH, METHOD_HEAD}, method); has {
		this.method = method
	} else {
		this.method = METHOD_GET
	}
	return this
}

func New() Request {
	instance := Request{
		initialed: true,
		header:    fasthttp.RequestHeader{},
	}
	return instance
}

func (this *Request) BaseUrl(url string) *Request {
	this.baseUrl = url
	return this
}

func (this *Request) Body(data []byte) *Request {
	this.payload = data

	return this
}

func (this *Request) JSON(obj interface{}) *Request {
	this.Header("Content-Type", "application/json")
	this.payload, this.error = json.Marshal(obj)

	return this
}

func (this *Request) Form(formData interface{}) *Request {
	switch reflect.TypeOf(formData).Kind() {
	case reflect.String:
		this.queryString = formData.(string)
	case reflect.Map:
	case reflect.Struct:
		this.queryString = buildQueryString(formData)
	}

	return this
}

func (this *Request) File(file interface{}) *Request {

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// 关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", file.(string))
	if err != nil {
		fmt.Println("error writing to buffer")
		return this
	}

	switch reflect.TypeOf(file).Kind() {
	case reflect.String:
		// 打开文件句柄操作
		fh, err := os.Open(file.(string))
		if err != nil {
			fmt.Println("error opening file")
			return this
		}
		defer fh.Close()
		// iocopy
		_, err = io.Copy(fileWriter, fh)
		if err != nil {
			fmt.Println("error copy file", err.Error())
			return this
		}
	case reflect.Array:
		fileWriter.Write(file.([]byte))
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	this.Header("Content-Type", contentType)
	this.payload = bodyBuf.Bytes()

	return this
}

func (this *Request) Query(query interface{}) *Request {
	switch reflect.TypeOf(query).Kind() {
	case reflect.String:
		this.queryString = query.(string)
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		this.queryString = buildQueryString(query)
	}

	return this
}

func (this *Request) Header(key string, value string) *Request {
	this.header.Set(key, value)

	return this
}

func (this *Request) Headers(headers fasthttp.RequestHeader) *Request {
	headers.CopyTo(&this.header)
	return this
}

func (this *Request) BearerToken(token string) *Request {

	this.Header("Authorization", token)

	return this
}

func (this *Request) Response() (response Response, err error) {
	if this.error != nil {
		return Response{}, this.error
	}

	url := this.baseUrl + this.url

	if len(this.queryString) > 0 {
		if strings.Contains(url, "?") {
			url += "&" + this.queryString
		} else {
			url += "?" + this.queryString
		}
	}

	httpRequest := fasthttp.AcquireRequest()
	httpResponse := fasthttp.AcquireResponse()
	httpRequest.SetRequestURI(url)
	httpRequest.SetBody(this.payload)
	this.header.CopyTo(&httpRequest.Header)
	httpRequest.Header.SetMethod(this.method)

	err = (&fasthttp.Client{}).Do(httpRequest, httpResponse)
	fasthttp.ReleaseRequest(httpRequest)
	if err != nil {
		return
	}

	err = response.From(httpResponse)
	fasthttp.ReleaseResponse(httpResponse)
	return
}

func buildQueryString(data interface{}) string {
	var slice []string

	if reflect.ValueOf(data).Kind() == reflect.Struct {
		data = structs.Map(data)
	}

	for k, v := range data.(map[string]interface{}) {

		switch reflect.TypeOf(v).Kind() {
		case reflect.Array:
			for _, sub := range v.([]interface{}) {
				slice = append(slice, k+"[]="+parseValue(sub))
			}
		default:
			slice = append(slice, k+"="+parseValue(v))
		}
	}
	return strings.Join(slice, "&")
}

func parseValue(value interface{}) string {
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return value.(string)
	case reflect.Int:
		return strconv.Itoa(value.(int))
	case reflect.Int64:
		return strconv.FormatInt(value.(int64), 10)
	case reflect.Uint32:
		return strconv.Itoa(int(value.(uint32)))
	case reflect.Bool:
		if value.(bool) {
			return "true"
		}
		return "false"
	case reflect.Float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	}

	return ""
}
