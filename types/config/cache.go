package config

type Cache struct {
	Enabled bool `yaml:"enabled" default:"false"`
	Dsn string `yaml:"dsn" default:"memory://heron"`
}
