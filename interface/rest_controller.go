package _interface

import (
	"github.com/gorilla/websocket"
	"github.com/duanchi/min/types"
	"github.com/gin-gonic/gin"
)

type RestControllerInterface interface {
	Fetch(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error)

	Create(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error)

	Update(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error)

	Remove(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error)

	Connect(connection *websocket.Conn, id string, resource string, parameters *gin.Params, ctx *gin.Context) (err types.Error)
}
