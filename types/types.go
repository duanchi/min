package types

import (
	"github.com/duanchi/min/server/httpserver/context"
)

type Route struct {
	Url     string
	Method  string
	Handler HandleFunc
}

type HandleFunc func(ctx *context.Context) error

type Response struct {
	RequestId string      `json:"request_id"`
	Status    bool        `json:"-"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

type PaginationData struct {
	Pagination Pagination  `json:"pagination"`
	Records    interface{} `json:"records"`
}

type Pagination struct {
	Total   int `json:"total"`
	Size    int `json:"size"`
	Pages   int `json:"pages"`
	Current int `json:"current"`
}

type Error interface {
	error
	Code() int
	Status() int
	Data() interface{}
}
