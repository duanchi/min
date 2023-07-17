package config

type HttpServer struct {
	Enabled    bool   `yaml:"enabled"`
	ServerHost string `yaml:"server-host" default:"${SERVER_HOST:}"`
	ServerPort string `yaml:"server-port" default:"${SERVER_PORT:9801}"`
	StaticPath string `yaml:"static-path" default:""`
}
