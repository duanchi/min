package types

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Context struct {
	gin.Context
}

type Param struct {
	gin.Param
}

type Params []Param

type WebsocketConnection struct {
	websocket.Conn
}
