package discovery

import (
	"github.com/duanchi/min/v2/abstract"
	"github.com/duanchi/min/v2/log"
	"github.com/duanchi/min/v2/microservice/discovery/nacos"
	"github.com/duanchi/min/v2/microservice/discovery/nacos/request"
	"github.com/duanchi/min/v2/microservice/discovery/nacos/response"
	"github.com/duanchi/min/v2/microservice/discovery/nacos/types"
	"github.com/duanchi/min/v2/types/config"
	"github.com/duanchi/min/v2/types/discovery"
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
		log.Log.Error("[min-framework]: Discovery Nacos get service list Error! " + err.Error())
		return serviceMap
	}

	if serviceList.Count > 0 {
		for _, serviceItem := range serviceList.ServiceList {
			service, err := this.GetService(serviceItem.Name)
			if err == nil {
				instances, listErr := this.GetAllInstances(serviceItem.Name)
				instanceMap := map[string][]discovery.Instance{}
				if listErr == nil {
					for _, serviceInstance := range instances {
						if _, has := instanceMap[serviceInstance.ClusterName]; !has {
							instanceMap[serviceInstance.ClusterName] = []discovery.Instance{}
						}
						instanceMap[serviceInstance.ClusterName] = append(instanceMap[serviceInstance.ClusterName], discovery.Instance{
							InstanceId:  serviceInstance.InstanceId,
							Ip:          serviceInstance.Ip,
							Port:        serviceInstance.Port,
							Weight:      serviceInstance.Weight,
							Healthy:     serviceInstance.Healthy,
							Enable:      serviceInstance.Enable,
							Ephemeral:   serviceInstance.Ephemeral,
							ServiceName: service.Name,
							Metadata:    serviceInstance.Metadata,
						})
					}
				}

				for clusterName, instances := range instanceMap {
					serviceMap[clusterName+"#"+serviceItem.GroupName+"@@"+serviceItem.Name] = discovery.Service{
						Instances: instances,
						Name:      serviceItem.Name,
						GroupName: serviceItem.GroupName,
					}
				}
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
func (this *NacosDiscovery) HeartBeat(heartBeat request.HeartBeat) {
	this.discoveryClient.HeartBeat(heartBeat)
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
