package route

import (
	"github.com/duanchi/min/config"
	"github.com/duanchi/min/rpc"
	"github.com/gin-gonic/gin"
)

func Init(httpServer *gin.Engine) {
	BaseRoutes.Init(httpServer)
	RestfulRoutesInit(httpServer)
	if config.Get("Rpc.Server.Enabled").(bool) == true {
		rpc.RpcBeans.Init(httpServer)
	}
}
