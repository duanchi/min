package abstract

import (
	"github.com/duanchi/min/v2/server/httpserver/context"
)

type Router struct {
	Bean
}

func (this *Router) Handle(path string, method string, params *context.Params, ctx *context.Context) {
}
