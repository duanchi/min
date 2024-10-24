package route

import (
	"github.com/duanchi/min/config"
	"github.com/duanchi/min/rpc"
	"github.com/duanchi/min/server/httpserver"
)

func Init(httpServer *httpserver.Httpserver) {
	BaseRouteInit(httpServer)
	RestfulRouteInit(httpServer)
	if config.Get("Rpc.Server.Enabled").(bool) == true {
		rpc.RpcBeans.Init(httpServer)
	}
}
