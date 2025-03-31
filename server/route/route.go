package route

import (
	"github.com/duanchi/min/v2/config"
	"github.com/duanchi/min/v2/rpc"
	"github.com/duanchi/min/v2/server/httpserver"
)

func Init(httpServer *httpserver.Httpserver) {
	BaseRouteInit(httpServer)
	RestfulRouteInit(httpServer)
	if config.Get("Rpc.Server.Enabled").(bool) == true {
		rpc.RpcBeans.Init(httpServer)
	}
}
