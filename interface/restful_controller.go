package _interface

import (
	"github.com/duanchi/min/v2/server/httpserver/context"
	"github.com/duanchi/min/v2/types"
)

type RestfulControllerInterface interface {
	Fetch(id string, resource string, parameters *context.Params, ctx *context.Context) (result interface{}, err types.Error)

	FetchList(id string, resource string, parameters *context.Params, ctx *context.Context) (result interface{}, err types.Error)

	Create(id string, resource string, parameters *context.Params, ctx *context.Context) (result interface{}, err types.Error)

	Update(id string, resource string, parameters *context.Params, ctx *context.Context) (result interface{}, err types.Error)

	Remove(id string, resource string, parameters *context.Params, ctx *context.Context) (result interface{}, err types.Error)

	Connect(connection *context.Websocket, id string, resource string, parameters *context.Params, ctx *context.Context) (err types.Error)
}
