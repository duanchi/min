package config

type Discovery struct {
	Enabled     bool   `yaml:"enabled" default:"false"`
	Url         string `yaml:"url" default:"nacos://127.0.0.1:8848/nacos/v1"`
	NamespaceId string `yaml:"namespace_id" default:""`
}
