package config

import "github.com/duanchi/min/types/config/feign"

type Feign struct {
	Enabled bool `yaml:"enabled" default:"false"`
	Services  []feign.Service `yaml:"services"`
	Debug    string `yaml:"debug" default:"false"`
}
