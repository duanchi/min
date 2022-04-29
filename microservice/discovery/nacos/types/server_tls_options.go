package types

var SkipVerifyConfig = TLSConfig{Enable: true}

func NewTLSConfig(opts ...TLSOption) *TLSConfig {
	tlsConfig := TLSConfig{Enable: true}
	for _, opt := range opts {
		opt(&tlsConfig)
	}
	return &tlsConfig
}

type TLSOption func(*TLSConfig)

func WithCA(caFile, serverNameOverride string) TLSOption {
	return func(tc *TLSConfig) {
		tc.CaFile = caFile
		tc.ServerNameOverride = serverNameOverride
	}
}

func WithCertificate(certFile, keyFile string) TLSOption {
	return func(tc *TLSConfig) {
		tc.CertFile = certFile
		tc.KeyFile = keyFile
	}
}
