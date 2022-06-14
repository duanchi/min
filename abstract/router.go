package abstract

import "github.com/duanchi/min/types"

type Router struct {
	Bean
}

func (this *Router) Handle(path string, method string, params types.Params, ctx *types.Context) {}
