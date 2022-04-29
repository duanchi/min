package config

type Discovery struct {
	Enabled        bool     `yaml:"enabled" default:"false"`
	Nodes          []string `yaml:"nodes"`
	NamespaceId    string   `yaml:"namespace_id" default:""`
	Group          string   `yaml:"group" default:"DEFAULT_GROUP"`
	UpdateInterval int64    `yaml:"update_interval" default:"10000"`
}
