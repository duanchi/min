package types

func NewServerConfig(ipAddr string, port uint64, opts ...ServerOption) *ServerConfig {
	serverConfig := &ServerConfig{
		IpAddr:      ipAddr,
		Port:        port,
		ContextPath: DEFAULT_CONTEXT_PATH,
		Scheme:      DEFAULT_SERVER_SCHEME,
	}

	for _, opt := range opts {
		opt(serverConfig)
	}

	return serverConfig
}

// ServerOption ...
type ServerOption func(*ServerConfig)

//WithScheme set Scheme for server
func WithScheme(scheme string) ServerOption {
	return func(config *ServerConfig) {
		config.Scheme = scheme
	}
}

//WithContextPath set contextPath for server
func WithContextPath(contextPath string) ServerOption {
	return func(config *ServerConfig) {
		config.ContextPath = contextPath
	}
}

//WithIpAddr set ip address for server
func WithIpAddr(ipAddr string) ServerOption {
	return func(config *ServerConfig) {
		config.IpAddr = ipAddr
	}
}

//WithPort set port for server
func WithPort(port string) ServerOption {
	return func(config *ServerConfig) {
		config.IpAddr = port
	}
}
