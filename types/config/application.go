package config

type Application struct {
	ServerPort string `yaml:"server_port" default:"${SERVER_PORT:9801}"`
	StaticPath string `yaml:"static_path" default:""`
}
