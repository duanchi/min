package request

import (
	"github.com/duanchi/min/v2/microservice/discovery/nacos/response"
)

type RegisterInstance struct {
	Ip          string            `param:"ip"`          //required
	Port        uint64            `param:"port"`        //required
	Weight      float64           `param:"weight"`      //required,it must be lager than 0
	Enable      bool              `param:"enabled"`     //required,the instance can be access or not
	Healthy     bool              `param:"healthy"`     //required,the instance is health or not
	Metadata    map[string]string `param:"metadata"`    //optional
	ClusterName string            `param:"clusterName"` //optional,default:DEFAULT
	ServiceName string            `param:"serviceName"` //required
	GroupName   string            `param:"groupName"`   //optional,default:DEFAULT_GROUP
	Ephemeral   bool              `param:"ephemeral"`   //optional
}

type DeregisterInstance struct {
	Ip          string `param:"ip"`          //required
	Port        uint64 `param:"port"`        //required
	Cluster     string `param:"cluster"`     //optional,default:DEFAULT
	ServiceName string `param:"serviceName"` //required
	GroupName   string `param:"groupName"`   //optional,default:DEFAULT_GROUP
	Ephemeral   bool   `param:"ephemeral"`   //optional
}

type HeartBeat struct {
	ServiceName string   `param:"serviceName"` //required
	GroupName   string   `param:"groupName"`   //optional,default:DEFAULT_GROUP
	Ephemeral   bool     `param:"ephemeral"`   //optional
	Beat        BeatInfo `param:"beat"`
	Healthy     bool     `param:"healthy"`
	Ip          string   `param:"ip"`
	Port        uint64   `param:"port"`
}

type BeatInfo struct {
	Ip          string            `json:"ip"`
	Port        uint64            `json:"port"`
	Weight      float64           `json:"weight"`
	ServiceName string            `json:"serviceName"`
	Cluster     string            `json:"cluster"`
	Metadata    map[string]string `json:"metadata"`
	Scheduled   bool              `json:"scheduled"`
}

type UpdateInstance struct {
	Ip          string            `param:"ip"`          //required
	Port        uint64            `param:"port"`        //required
	Weight      float64           `param:"weight"`      //required,it must be lager than 0
	Enable      bool              `param:"enabled"`     //required,the instance can be access or not
	Healthy     bool              `param:"healthy"`     //required,the instance is health or not
	Metadata    map[string]string `param:"metadata"`    //optional
	ClusterName string            `param:"clusterName"` //optional,default:DEFAULT
	ServiceName string            `param:"serviceName"` //required
	GroupName   string            `param:"groupName"`   //optional,default:DEFAULT_GROUP
	Ephemeral   bool              `param:"ephemeral"`   //optional
}

type GetService struct {
	Clusters    []string `param:"clusters"`    //optional,default:DEFAULT
	ServiceName string   `param:"serviceName"` //required
	GroupName   string   `param:"groupName"`   //optional,default:DEFAULT_GROUP
}

type GetAllServiceInfo struct {
	NameSpace string `param:"nameSpace"` //optional,default:public
	GroupName string `param:"groupName"` //optional,default:DEFAULT_GROUP
	PageNo    uint32 `param:"pageNo"`    //optional,default:1
	PageSize  uint32 `param:"pageSize"`  //optional,default:10
}

type Subscribe struct {
	ServiceName       string                                        `param:"serviceName"` //required
	Clusters          []string                                      `param:"clusters"`    //optional,default:DEFAULT
	GroupName         string                                        `param:"groupName"`   //optional,default:DEFAULT_GROUP
	SubscribeCallback func(services []response.Instance, err error) //required
}

type SelectAllInstances struct {
	Clusters    []string `param:"clusters"`    //optional,default:DEFAULT
	ServiceName string   `param:"serviceName"` //required
	GroupName   string   `param:"groupName"`   //optional,default:DEFAULT_GROUP
}

type SelectInstances struct {
	Clusters    []string `param:"clusters"`    //optional,default:DEFAULT
	ServiceName string   `param:"serviceName"` //required
	GroupName   string   `param:"groupName"`   //optional,default:DEFAULT_GROUP
	HealthyOnly bool     `param:"healthyOnly"` //optional,return only healthy instance
}

type SelectOneHealthInstance struct {
	Clusters    []string `param:"clusters"`    //optional,default:DEFAULT
	ServiceName string   `param:"serviceName"` //required
	GroupName   string   `param:"groupName"`   //optional,default:DEFAULT_GROUP
}
