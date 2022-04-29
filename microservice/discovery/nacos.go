package discovery

import (
	"fmt"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/types/config"
	"github.com/duanchi/min/types/discovery"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosDiscovery struct {
	abstract.Bean
	serverConfig      []constant.ServerConfig
	clientConfig      constant.ClientConfig
	applicationConfig config.Application
	discoveryConfig   config.Discovery
	namingClient      naming_client.INamingClient
}

func (this *NacosDiscovery) Init() {
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &this.clientConfig,
			ServerConfigs: this.serverConfig,
		},
	)
	if err != nil {
		fmt.Println("[min-framework]: Discovery Create Nacos Naming Client Error! " + err.Error())
	} else {
		fmt.Println("[min-framework]: Discovery Connect Nacos Naming Client Success!")
	}
	this.namingClient = namingClient

	ServiceMap = this.GetServiceList()
}

func (this *NacosDiscovery) GetServiceList() map[string]discovery.Service {
	serviceMap := map[string]discovery.Service{}
	serviceList, err := this.namingClient.GetAllServicesInfo(
		vo.GetAllServiceInfoParam{
			NameSpace: this.discoveryConfig.NamespaceId,
			PageNo:    1,
			PageSize:  512,
		})
	if err != nil {
		fmt.Println("[min-framework]: Discovery Nacos get service list Error! " + err.Error())
		return serviceMap
	}

	if serviceList.Count > 0 {
		for _, serviceName := range serviceList.Doms {
			this.GetAllInstances(serviceName, this.discoveryConfig.Group)
			service, err := this.GetService(serviceName, this.discoveryConfig.Group)

			if err == nil {
				serviceMap[serviceName] = service
			}
		}
	}

	return serviceMap
}

func (this *NacosDiscovery) RegisterInstance()   {}
func (this *NacosDiscovery) DeregisterInstance() {}
func (this *NacosDiscovery) GetService(name string, group string) (discoveryService discovery.Service, err error) {
	service, err := this.namingClient.GetService(vo.GetServiceParam{
		ServiceName: name,
		Clusters:    []string{"DEFAULT"}, // default value is DEFAULT
		GroupName:   group,               // default value is DEFAULT_GROUP
	})

	discoveryService.Name = service.Name
	discoveryService.GroupName = service.GroupName
	var instances []discovery.Instance
	for _, instance := range service.Hosts {
		instances = append(instances, discovery.Instance{
			InstanceId:        instance.InstanceId,
			Ip:                instance.Ip,
			Port:              instance.Port,
			Weight:            instance.Weight,
			Healthy:           instance.Healthy,
			Enable:            instance.Enable,
			Ephemeral:         instance.Ephemeral,
			ServiceName:       service.Name,
			Metadata:          instance.Metadata,
			HeartBeatInterval: instance.InstanceHeartBeatInterval,
			HeartBeatTimeOut:  instance.InstanceHeartBeatTimeOut,
		})
	}

	discoveryService.Instances = instances
	return
}
func (this *NacosDiscovery) GetAllInstances(serviceName string, group string) {
	instances, err := this.namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: serviceName,
		Clusters:    []string{"DEFAULT"}, // default value is DEFAULT
		GroupName:   group,               // default value is DEFAULT_GROUP
	})

	fmt.Println(instances, err)
}
func (this *NacosDiscovery) GetInstances(serviceName string, group string) {
	instances, err := this.namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		Clusters:    []string{"DEFAULT"}, // default value is DEFAULT
		GroupName:   group,               // default value is DEFAULT_GROUP
		HealthyOnly: true,
	})

	fmt.Println(instances, err)
}
func (this *NacosDiscovery) GetHealthInstance() {}
func (this *NacosDiscovery) Subscribe() {

}
func (this *NacosDiscovery) UnSubscribe() {}
