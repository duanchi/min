package types

import (
	"time"
)

func NewClientConfig(opts ...ClientOption) *ClientConfig {
	clientConfig := &ClientConfig{
		TimeoutMs:            10 * 1000,
		BeatInterval:         5 * 1000,
		OpenKMS:              false,
		UpdateThreadNum:      20,
		NotLoadCacheAtStart:  false,
		UpdateCacheWhenEmpty: false,
		LogLevel:             "info",
	}

	for _, opt := range opts {
		opt(clientConfig)
	}

	return clientConfig
}

// ClientOption ...
type ClientOption func(*ClientConfig)

// WithTimeoutMs ...
func WithTimeoutMs(timeoutMs uint64) ClientOption {
	return func(config *ClientConfig) {
		config.TimeoutMs = timeoutMs
	}
}

// WithAppName ...
func WithAppName(appName string) ClientOption {
	return func(config *ClientConfig) {
		config.AppName = appName
	}
}

// WithBeatInterval ...
func WithBeatInterval(beatInterval int64) ClientOption {
	return func(config *ClientConfig) {
		config.BeatInterval = beatInterval
	}
}

// WithNamespaceId ...
func WithNamespaceId(namespaceId string) ClientOption {
	return func(config *ClientConfig) {
		config.NamespaceId = namespaceId
	}
}

// WithEndpoint ...
func WithEndpoint(endpoint string) ClientOption {
	return func(config *ClientConfig) {
		config.Endpoint = endpoint
	}
}

// WithRegionId ...
func WithRegionId(regionId string) ClientOption {
	return func(config *ClientConfig) {
		config.RegionId = regionId
	}
}

// WithAccessKey ...
func WithAccessKey(accessKey string) ClientOption {
	return func(config *ClientConfig) {
		config.AccessKey = accessKey
	}
}

// WithSecretKey ...
func WithSecretKey(secretKey string) ClientOption {
	return func(config *ClientConfig) {
		config.SecretKey = secretKey
	}
}

// WithOpenKMS ...
func WithOpenKMS(openKMS bool) ClientOption {
	return func(config *ClientConfig) {
		config.OpenKMS = openKMS
	}
}

// WithUpdateThreadNum ...
func WithUpdateThreadNum(updateThreadNum int) ClientOption {
	return func(config *ClientConfig) {
		config.UpdateThreadNum = updateThreadNum
	}
}

// WithNotLoadCacheAtStart ...
func WithNotLoadCacheAtStart(notLoadCacheAtStart bool) ClientOption {
	return func(config *ClientConfig) {
		config.NotLoadCacheAtStart = notLoadCacheAtStart
	}
}

// WithUpdateCacheWhenEmpty ...
func WithUpdateCacheWhenEmpty(updateCacheWhenEmpty bool) ClientOption {
	return func(config *ClientConfig) {
		config.UpdateCacheWhenEmpty = updateCacheWhenEmpty
	}
}

// WithUsername ...
func WithUsername(username string) ClientOption {
	return func(config *ClientConfig) {
		config.Username = username
	}
}

// WithPassword ...
func WithPassword(password string) ClientOption {
	return func(config *ClientConfig) {
		config.Password = password
	}
}

// WithLogLevel ...
func WithLogLevel(logLevel string) ClientOption {
	return func(config *ClientConfig) {
		config.LogLevel = logLevel
	}
}

// WithLogSampling ...
func WithLogSampling(tick time.Duration, initial int, thereafter int) ClientOption {
	return func(config *ClientConfig) {
		config.LogSampling = &ClientLogSamplingConfig{initial, thereafter, tick}
	}
}

// WithLogRollingConfig ...
func WithLogRollingConfig(rollingConfig *ClientLogRollingConfig) ClientOption {
	return func(config *ClientConfig) {
		config.LogRollingConfig = rollingConfig
	}
}

func WithTLS(tlsCfg TLSConfig) ClientOption {
	return func(config *ClientConfig) {
		config.TLSCfg = tlsCfg
	}
}
