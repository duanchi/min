package abstract

import (
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/duanchi/min/types"
	"github.com/gorilla/websocket"
)

type RestfulController struct {
	Bean
}

func (this *RestfulController) Fetch(id string, resource string, parameters *context.Params, ctx *context.Context) (result interface{}, err types.Error) {
	return "error", nil
}

func (this *RestfulController) FetchList(id string, resource string, parameters *context.Params, ctx *context.Context) (result interface{}, err types.Error) {
	return "error", nil
}

func (this *RestfulController) Create(id string, resource string, parameters *context.Params, ctx *context.Context) (result interface{}, err types.Error) {
	return "error", nil
}

func (this *RestfulController) Update(id string, resource string, parameters *context.Params, ctx *context.Context) (result interface{}, err types.Error) {
	return "error", nil
}

func (this *RestfulController) Remove(id string, resource string, parameters *context.Params, ctx *context.Context) (result interface{}, err types.Error) {
	return "error", nil
}

func (this *RestfulController) Connect(connection *websocket.Conn, id string, resource string, parameters *context.Params, ctx *context.Context) (err types.Error) {
	return nil
}
