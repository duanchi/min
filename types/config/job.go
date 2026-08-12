package config

type Job struct {
	Enabled        bool   `yaml:"enabled" default:"false"`
	Server         string `yaml:"server"`
	AccessToken    string `yaml:"access-token"`
	ExecutorName   string `yaml:"executor-name"`
	EnableRegister bool   `yaml:"enable-register" default:"false"`
	ExecutorIp     string `yaml:"executor-ip"`
	ExecutorPort   string `yaml:"executor-port"`
}
