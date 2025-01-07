package config

type HttpServer struct {
	Enabled    bool          `yaml:"enabled"`
	ServerHost string        `yaml:"server-host" default:"${SERVER_HOST:}"`
	ServerPort string        `yaml:"server-port" default:"${SERVER_PORT:9801}"`
	StaticPath string        `yaml:"static-path" default:""`
	Restful    RestfulConfig `yaml:"restful"`
}

type RestfulConfig struct {
	CustomResponse bool `yaml:"custom-response"`
}
