package discovery

type Service struct {
	Instances []Instance `json:"instances"`
	Name      string     `json:"name"`
	GroupName string     `json:"groupName"`
}
