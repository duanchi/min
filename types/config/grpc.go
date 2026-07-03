package config

type Grpc struct {
	Server GrpcServer `yaml:"server"`
}

type GrpcServer struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}
