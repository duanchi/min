package route

import (
	"github.com/duanchi/min/config"
	"github.com/duanchi/min/rpc"
	"github.com/gofiber/fiber/v2"
)

func Init(httpServer *fiber.App) {
	BaseRoutes.Init(httpServer)
	RestfulRoutesInit(httpServer)
	if config.Get("Rpc.Server.Enabled").(bool) == true {
		rpc.RpcBeans.Init(httpServer)
	}
}
