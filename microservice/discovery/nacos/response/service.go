package response

import "time"

const (
	StateRunning = iota
	StateShutdown
)

type Instance struct {
	InstanceId                string            `json:"instanceId"`
	Ip                        string            `json:"ip"`
	Port                      uint64            `json:"port"`
	Weight                    float64           `json:"weight"`
	Healthy                   bool              `json:"healthy"`
	Enable                    bool              `json:"enabled"`
	Ephemeral                 bool              `json:"ephemeral"`
	ClusterName               string            `json:"clusterName"`
	ServiceName               string            `json:"serviceName"`
	Metadata                  map[string]string `json:"metadata"`
	InstanceHeartBeatInterval int               `json:"instanceHeartBeatInterval"`
	IpDeleteTimeout           int               `json:"ipDeleteTimeout"`
	InstanceHeartBeatTimeOut  int               `json:"instanceHeartBeatTimeOut"`
}

type Service struct {
	CacheMillis              uint64     `json:"cacheMillis"`
	Hosts                    []Instance `json:"hosts"`
	Checksum                 string     `json:"checksum"`
	LastRefTime              uint64     `json:"lastRefTime"`
	Clusters                 []Cluster  `json:"clusters"`
	Name                     string     `json:"name"`
	GroupName                string     `json:"groupName"`
	Valid                    bool       `json:"valid"`
	AllIPs                   bool       `json:"allIPs"`
	ReachProtectionThreshold bool       `json:"reachProtectionThreshold"`
}

type InstanceResult struct {
	CacheMillis              uint64     `json:"cacheMillis"`
	Hosts                    []Instance `json:"hosts"`
	Checksum                 string     `json:"checksum"`
	LastRefTime              uint64     `json:"lastRefTime"`
	Name                     string     `json:"name"`
	GroupName                string     `json:"groupName"`
	Valid                    bool       `json:"valid"`
	AllIPs                   bool       `json:"allIPs"`
	ReachProtectionThreshold bool       `json:"reachProtectionThreshold"`
}

type ServiceDetail struct {
	Service  ServiceInfo `json:"service"`
	Clusters []Cluster   `json:"clusters"`
}

type ServiceInfo struct {
	App              string            `json:"app"`
	Group            string            `json:"group"`
	HealthCheckMode  string            `json:"healthCheckMode"`
	Metadata         map[string]string `json:"metadata"`
	Name             string            `json:"name"`
	ProtectThreshold float64           `json:"protectThreshold"`
	Selector         ServiceSelector   `json:"selector"`
}

type ServiceSelector struct {
	Selector string
}

type Cluster struct {
	ServiceName      string               `json:"serviceName"`
	Name             string               `json:"name"`
	HealthyChecker   ClusterHealthChecker `json:"healthyChecker"`
	DefaultPort      uint64               `json:"defaultPort"`
	DefaultCheckPort uint64               `json:"defaultCheckPort"`
	UseIPPort4Check  bool                 `json:"useIpPort4Check"`
	Metadata         map[string]string    `json:"metadata"`
}

type ClusterHealthChecker struct {
	Type string `json:"type"`
}

type BeatInfo struct {
	Ip          string            `json:"ip"`
	Port        uint64            `json:"port"`
	Weight      float64           `json:"weight"`
	ServiceName string            `json:"serviceName"`
	Cluster     string            `json:"cluster"`
	Metadata    map[string]string `json:"metadata"`
	Scheduled   bool              `json:"scheduled"`
	Period      time.Duration     `json:"-"`
	State       int32             `json:"-"`
}

type ExpressionSelector struct {
	Type       string `json:"type"`
	Expression string `json:"expression"`
}

/*type ServiceList struct {
	Count int64    `json:"count"`
	Doms  []string `json:"doms"`
}*/

type ServiceList struct {
	Count       int64         `json:"count"`
	ServiceList []ServiceItem `json:"serviceList"`
}

type ServiceItem struct {
	Name                 string `json:"name"`
	GroupName            string `json:"groupName"`
	ClusterCount         int    `json:"clusterCount"`
	IpCount              int    `json:"ipCount"`
	HealthyInstanceCount int    `json:"healthyInstanceCount"`
	TriggerFlag          string `json:"triggerFlag"`
}
