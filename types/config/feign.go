package config

import "github.com/duanchi/min/v2/types/config/feign"

type Feign struct {
	Enabled bool          `yaml:"enabled" default:"false"`
	Service feign.Service `yaml:"service"`
	Client  feign.Client  `yaml:"client"`
	Debug   string        `yaml:"debug" default:"false"`
}
