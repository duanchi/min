package request

import (
	"github.com/duanchi/min/v2/microservice/discovery/nacos/types"
	"github.com/duanchi/min/v2/types/config"
)

type Client struct {
	ClientConfig  *types.ClientConfig  // optional
	ServerConfigs []types.ServerConfig // optional
	RuntimeConfig config.Discovery
}
