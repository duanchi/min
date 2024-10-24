package _interface

import (
	"github.com/duanchi/min/server/httpserver/context"
)

type RouterInterface interface {
	Handle(path string, method string, params *context.Params, ctx *context.Context) error
}
