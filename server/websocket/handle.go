package websocket

import (
	"github.com/duanchi/min/server/httpserver/context"
	"github.com/duanchi/min/types"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handle(id string, resource string, parameters *context.Params, ctx *context.Context,
	handleFunction func(connection *websocket.Conn, id string, resource string, parameters *context.Params, ctx *context.Context) (err types.Error),
) (err error) {

	/*connection, websocketError := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if websocketError != nil {
		return types.RuntimeError{
			Message:   websocketError.Error(),
			ErrorCode: http.StatusInternalServerError,
		}
	}
	defer connection.Close()
	err = handleFunction(connection, id, resource, parameters, ctx)
	if err != nil {
		return
	}*/

	return
}
