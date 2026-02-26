package config

type HttpServer struct {
	Enabled    bool          `yaml:"enabled"`
	ServerHost string        `yaml:"server-host" default:"${SERVER_HOST:}"`
	ServerPort string        `yaml:"server-port" default:"${SERVER_PORT:9801}"`
	StaticPath string        `yaml:"static-path" default:""`
	Restful    RestfulConfig `yaml:"restful"`
	Config     ServerConfig  `yaml:"config"`
}

type ServerConfig struct {
	AppName           string `yaml:"app-name" default:"SERVER_CONFIG_APP_NAME:"`
	ClientMaxBodySize string `yaml:"client-max-body-size" default:"${SERVER_CONFIG_CLIENT_MAX_BODY_SIZE:50M}"`
	Concurrency       int    `yaml:"concurrency" default:"${SERVER_CONFIG_CONCURRENCY:262144}"`
}

type RestfulConfig struct {
	CustomResponse bool `yaml:"custom-response"`
}
