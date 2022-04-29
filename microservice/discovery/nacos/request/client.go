package request

import (
	"github.com/duanchi/min/microservice/discovery/nacos/types"
	"github.com/duanchi/min/types/config"
)

type Client struct {
	ClientConfig  *types.ClientConfig  // optional
	ServerConfigs []types.ServerConfig // optional
	RuntimeConfig config.Discovery
}
