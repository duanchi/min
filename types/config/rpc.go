package config

import "github.com/duanchi/min/types/config/rpc"

type Rpc struct {
	Server rpc.Server `yaml:"server"`
}
