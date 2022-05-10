package discovery

import (
	"fmt"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/microservice/discovery/nacos"
	"github.com/duanchi/min/microservice/discovery/nacos/request"
	"github.com/duanchi/min/microservice/discovery/nacos/response"
	"github.com/duanchi/min/microservice/discovery/nacos/types"
	"github.com/duanchi/min/types/config"
	"github.com/duanchi/min/types/discovery"
)

type NacosDiscovery struct {
	abstract.Bean
	serverConfig      []types.ServerConfig
	clientConfig      types.ClientConfig
	applicationConfig config.Application
	discoveryConfig   config.Discovery
	discoveryClient   nacos.DiscoveryClient
}

func (this *NacosDiscovery) Init() {

	this.discoveryClient = nacos.NewDiscoveryClient(request.Client{
		ClientConfig:  &this.clientConfig,
		ServerConfigs: this.serverConfig,
		RuntimeConfig: this.discoveryConfig,
	})

	ServiceMap = this.GetServiceList()
}

func (this *NacosDiscovery) GetServiceList() map[string]discovery.Service {
	serviceMap := map[string]discovery.Service{}
	serviceList, err := this.discoveryClient.GetAllServicesInfo()
	if err != nil {
		fmt.Println("[min-framework]: Discovery Nacos get service list Error! " + err.Error())
		return serviceMap
	}

	if serviceList.Count > 0 {
		for _, serviceName := range serviceList.Doms {
			service, err := this.GetService(serviceName)
			if err == nil {
				instances, listErr := this.GetAllInstances(serviceName)
				if listErr == nil {
					for _, instance := range instances {
						service.Instances = append(service.Instances, discovery.Instance{
							InstanceId:  instance.InstanceId,
							Ip:          instance.Ip,
							Port:        instance.Port,
							Weight:      instance.Weight,
							Healthy:     instance.Healthy,
							Enable:      instance.Enable,
							Ephemeral:   instance.Ephemeral,
							ServiceName: service.Name,
							Metadata:    instance.Metadata,
						})
					}
				}
				serviceMap[serviceName] = service
			}
		}
	}

	return serviceMap
}

func (this *NacosDiscovery) RegisterInstance(instance request.RegisterInstance) {
	this.discoveryClient.RegisterInstance(instance)
}
func (this *NacosDiscovery) DeregisterInstance(instance request.DeregisterInstance) {
	this.discoveryClient.DeregisterInstance(instance)
}
func (this *NacosDiscovery) GetService(serviceName string) (discoveryService discovery.Service, err error) {
	service, err := this.discoveryClient.GetService(serviceName)

	discoveryService.Name = service.Name
	discoveryService.GroupName = service.GroupName
	var instances []discovery.Instance
	for _, instance := range service.Hosts {
		instances = append(instances, discovery.Instance{
			InstanceId:  instance.InstanceId,
			Ip:          instance.Ip,
			Port:        instance.Port,
			Weight:      instance.Weight,
			Healthy:     instance.Healthy,
			Enable:      instance.Enable,
			Ephemeral:   instance.Ephemeral,
			ServiceName: service.Name,
			Metadata:    instance.Metadata,
		})
	}

	discoveryService.Instances = instances
	return
}
func (this *NacosDiscovery) GetAllInstances(serviceName string) (instances []response.Instance, err error) {
	instances, _ = this.discoveryClient.SelectAllInstances(serviceName)
	return
}
func (this *NacosDiscovery) GetInstances(serviceName string) (instances []response.Instance, err error) {
	instances, _ = this.discoveryClient.SelectInstances(serviceName)
	return
}
