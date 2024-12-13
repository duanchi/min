package websocket

import (
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/duanchi/min/types"
	"github.com/gofiber/contrib/websocket"
)

func Handle(id string, resource string, parameters *context.Params, ctx *context.Context,
	handleFunction func(connection *context.Websocket, id string, resource string, parameters *context.Params, ctx *context.Context) (err types.Error),
) (err error) {

	websocketHandle := websocket.New(func(conn *websocket.Conn) {
		defer conn.Close()
		connection := context.NewWebsocket(conn)
		for {
			err = handleFunction(connection, id, resource, parameters, ctx)
			if err != nil {
				return
			}
		}
	})

	err = websocketHandle(ctx.Ctx())
	if err != nil {
		return
	}
	return
}
