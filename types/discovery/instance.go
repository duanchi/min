package discovery

type Instance struct {
	InstanceId        string            `json:"instanceId"`
	Ip                string            `json:"ip"`
	Port              uint64            `json:"port"`
	Weight            float64           `json:"weight"`
	Healthy           bool              `json:"healthy"`
	Enable            bool              `json:"enabled"`
	Ephemeral         bool              `json:"ephemeral"`
	ServiceName       string            `json:"serviceName"`
	Metadata          map[string]string `json:"metadata"`
	HeartBeatInterval int               `json:"heartBeatInterval"`
	HeartBeatTimeOut  int               `json:"heartBeatTimeOut"`
}
