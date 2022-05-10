package nacos

import (
	"encoding/json"
	"github.com/duanchi/min/microservice/discovery/nacos/holder"
	"github.com/duanchi/min/microservice/discovery/nacos/request"
	"github.com/duanchi/min/microservice/discovery/nacos/response"
)

type DiscoveryClient struct {
	config        request.Client
	requestHolder *holder.HttpHolder
}

func NewDiscoveryClient(discoveryConfig request.Client) DiscoveryClient {
	client := DiscoveryClient{
		config:        discoveryConfig,
		requestHolder: holder.NewHttpHolder(discoveryConfig),
	}
	return client
}

func (this *DiscoveryClient) RegisterInstance(param request.RegisterInstance) (ok bool, err error) {
	metadataString, _ := json.Marshal(param.Metadata)
	ok, err = this.requestHolder.POST("/ns/instance", map[string]interface{}{
		"ip":          param.Ip,
		"port":        param.Port,
		"weight":      param.Weight,
		"enable":      param.Enable,
		"healthy":     param.Healthy,
		"metadata":    string(metadataString),
		"clusterName": "DEFAULT",
		"serviceName": param.ServiceName,
		"groupName":   param.GroupName,
		"ephemeral":   param.Ephemeral,
	})
	return
}

func (this *DiscoveryClient) DeregisterInstance(param request.DeregisterInstance) (ok bool, err error) {
	ok, err = this.requestHolder.DELETE("/ns/instance", map[string]interface{}{
		"ip":          param.Ip,
		"port":        param.Port,
		"clusterName": "DEFAULT",
		"serviceName": param.ServiceName,
		"groupName":   param.GroupName,
		"ephemeral":   param.Ephemeral,
	})
	return
}

func (this *DiscoveryClient) UpdateInstance(param request.UpdateInstance) (ok bool, err error) {
	metadataString, _ := json.Marshal(param.Metadata)
	ok, err = this.requestHolder.PUT("/ns/instance", map[string]interface{}{
		"ip":          param.Ip,
		"port":        param.Port,
		"weight":      param.Weight,
		"enable":      param.Enable,
		"healthy":     param.Healthy,
		"metadata":    string(metadataString),
		"clusterName": "DEFAULT",
		"serviceName": param.ServiceName,
		"groupName":   param.GroupName,
		"ephemeral":   param.Ephemeral,
	})
	return
}

func (this *DiscoveryClient) GetService(serviceName string) (response response.Service, err error) {
	err = this.requestHolder.GET("/ns/service", map[string]interface{}{
		"serviceName": serviceName,
		"namespaceId": this.config.ClientConfig.NamespaceId,
		"groupName":   this.config.RuntimeConfig.Group,
	}, &response)
	return
}

func (this *DiscoveryClient) SelectAllInstances(serviceName string) (instanceResponse []response.Instance, err error) {
	serviceResponse := response.InstanceResult{}
	err = this.requestHolder.GET("/ns/instance/list", map[string]interface{}{
		"serviceName": serviceName,
		"namespaceId": this.config.ClientConfig.NamespaceId,
		"groupName":   this.config.RuntimeConfig.Group,
		"cluster":     "DEFAULT",
	}, &serviceResponse)
	if err == nil {
		instanceResponse = serviceResponse.Hosts
	}
	return
}

func (this *DiscoveryClient) SelectInstances(serviceName string) (instanceResponse []response.Instance, err error) {
	serviceResponse := response.InstanceResult{}
	err = this.requestHolder.GET("/ns/instance/list", map[string]interface{}{
		"serviceName": serviceName,
		"namespaceId": this.config.ClientConfig.NamespaceId,
		"groupName":   this.config.RuntimeConfig.Group,
		"cluster":     "DEFAULT",
		"healthyOnly": true,
	}, &serviceResponse)
	if err == nil {
		instanceResponse = serviceResponse.Hosts
	}
	return
}

func (this *DiscoveryClient) GetAllServicesInfo() (response response.ServiceList, err error) {
	err = this.requestHolder.GET("/ns/service/list", map[string]interface{}{
		"pageNo":      1,
		"pageSize":    512,
		"namespaceId": this.config.ClientConfig.NamespaceId,
		"groupName":   this.config.RuntimeConfig.Group,
	}, &response)
	return
}
