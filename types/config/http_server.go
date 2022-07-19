package config

type HttpServer struct {
	Enabled    bool   `yaml:"enabled"`
	ServerPort string `yaml:"server-port" default:"${SERVER_PORT:9801}"`
	StaticPath string `yaml:"static-path" default:""`
}
