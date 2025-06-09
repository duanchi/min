package discovery

type Client struct {
	Enabled           bool              `yaml:"enabled" default:"false"`
	Ip                string            `yaml:"ip"`
	Port              string            `yaml:"port"`
	Scheme            string            `yaml:"scheme" default:"http"`
	InstanceId        string            `yaml:"instance_id"`
	Metadata          map[string]string `yaml:"metadata"`
	HeartbeatInterval int64             `yaml:"heartbeat_interval" default:"4000"`
}
