package rpc

type Server struct {
	Enabled bool `yaml:"enabled" value:"false"`
	Prefix string `yaml:"prefix" value:""`
}
