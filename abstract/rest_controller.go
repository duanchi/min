package abstract

import (
	"github.com/duanchi/min/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type RestController struct {
	Bean
}

func (this *RestController) Fetch(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	return "error", nil
}

func (this *RestController) FetchList(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	return "error", nil
}

func (this *RestController) Create(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	return "error", nil
}

func (this *RestController) Update(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	return "error", nil
}

func (this *RestController) Remove(id string, resource string, parameters *gin.Params, ctx *gin.Context) (result interface{}, err types.Error) {
	return "error", nil
}

func (this *RestController) Connect(connection *websocket.Conn, id string, resource string, parameters *gin.Params, ctx *gin.Context) (err types.Error) {
	return nil
}
