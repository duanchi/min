package config

import "github.com/duanchi/min/v2/types/config/rpc"

type Rpc struct {
	Server rpc.Server `yaml:"server"`
}
