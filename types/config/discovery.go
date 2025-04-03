package config

import "github.com/duanchi/min/v2/types/config/discovery"

type Discovery struct {
	Enabled        bool             `yaml:"enabled" default:"false"`
	Nodes          string           `yaml:"nodes"`
	NamespaceId    string           `yaml:"namespace_id" default:""`
	Group          string           `yaml:"group" default:""`
	UpdateInterval int64            `yaml:"update_interval" default:"10000"`
	Client         discovery.Client `yaml:"client"`
	Weight         float64          `yaml:"weight" default:"10"`
}
