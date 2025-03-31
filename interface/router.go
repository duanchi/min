package _interface

import (
	"github.com/duanchi/min/v2/server/httpserver/context"
)

type RouterInterface interface {
	Handle(path string, method string, params *context.Params, ctx *context.Context)
}
