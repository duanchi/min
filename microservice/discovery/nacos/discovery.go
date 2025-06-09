package nacos

import (
	"encoding/json"
	"github.com/duanchi/min/v2/microservice/discovery/nacos/holder"
	"github.com/duanchi/min/v2/microservice/discovery/nacos/request"
	"github.com/duanchi/min/v2/microservice/discovery/nacos/response"
	"net/url"
)

type DiscoveryClient struct {
	config request.Client
	// requestHolder *holder.HttpHolder
}

func NewDiscoveryClient(discoveryConfig request.Client) DiscoveryClient {
	client := DiscoveryClient{
		config: discoveryConfig,
		// requestHolder: holder.NewHttpHolder(discoveryConfig),
	}
	return client
}

func (this *DiscoveryClient) getRequestHolder() *holder.HttpHolder {
	return holder.NewHttpHolder(this.config)
}

func (this *DiscoveryClient) RegisterInstance(param request.RegisterInstance) (ok bool, err error) {
	metadataString, _ := json.Marshal(param.Metadata)
	ok, err = this.getRequestHolder().POST(this.parseUrl("/ns/instance"), map[string]interface{}{
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
	ok, err = this.getRequestHolder().DELETE(this.parseUrl("/ns/instance"), map[string]interface{}{
		"ip":          param.Ip,
		"port":        param.Port,
		"clusterName": "DEFAULT",
		"serviceName": param.ServiceName,
		"groupName":   param.GroupName,
		"ephemeral":   param.Ephemeral,
	})
	return
}

func (this *DiscoveryClient) parseUrl(path string, version ...string) string {
	if len(version) > 0 {
		return "/nacos/" + version[0] + path
	} else {
		return "/nacos/v2" + path
	}
}

func (this *DiscoveryClient) HeartBeat(param request.HeartBeat) (ok bool, err error) {
	beatBytes, err := json.Marshal(param.Beat)
	_, err = this.getRequestHolder().
		Holder().
		Url(this.parseUrl("/ns/instance/beat", "v1")).
		Method("PUT").
		Query(map[string]interface{}{
			"serviceName": param.ServiceName,
			"groupName":   param.GroupName,
			"ip":          param.Ip,
			"port":        param.Port,
			"healthy":     param.Healthy,
			"ephemeral":   param.Ephemeral,
			"beat":        url.QueryEscape(string(beatBytes)),
		}).
		Response()
	/*ok, err = this.getRequestHolder().PUT("/ns/instance/beat", map[string]interface{}{
		"serviceName": param.ServiceName,
		"groupName":   param.GroupName,
		"ip":          param.Ip,
		"port":        param.Port,
		"healthy":     param.Healthy,
		"ephemeral":   param.Ephemeral,
		"beat":        url.QueryEscape(string(beatBytes)),
	})*/
	return err == nil, err
}

func (this *DiscoveryClient) UpdateInstance(param request.UpdateInstance) (ok bool, err error) {
	metadataString, _ := json.Marshal(param.Metadata)
	ok, err = this.getRequestHolder().PUT(this.parseUrl("/ns/instance"), map[string]interface{}{
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

func (this *DiscoveryClient) GetService(serviceName string) (service response.Service, err error) {
	res := response.Result[response.Service]{}
	err = this.getRequestHolder().GET(this.parseUrl("/ns/service"), map[string]interface{}{
		"serviceName": serviceName,
		"namespaceId": this.config.ClientConfig.NamespaceId,
		"groupName":   this.config.RuntimeConfig.Group,
	}, &res)
	return res.Data, err
}

func (this *DiscoveryClient) SelectAllInstances(serviceName string) (instanceResponse []response.Instance, err error) {
	serviceResponse := response.Result[response.InstanceResult]{}
	err = this.getRequestHolder().GET(this.parseUrl("/ns/instance/list"), map[string]interface{}{
		"serviceName": serviceName,
		"namespaceId": this.config.ClientConfig.NamespaceId,
		"groupName":   this.config.RuntimeConfig.Group,
		"cluster":     "DEFAULT",
	}, &serviceResponse)
	if err == nil {
		instanceResponse = serviceResponse.Data.Hosts
	}
	return
}

func (this *DiscoveryClient) SelectInstances(serviceName string) (instanceResponse []response.Instance, err error) {
	serviceResponse := response.Result[response.InstanceResult]{}
	err = this.getRequestHolder().GET(this.parseUrl("/ns/instance/list"), map[string]interface{}{
		"serviceName": serviceName,
		"namespaceId": this.config.ClientConfig.NamespaceId,
		"groupName":   this.config.RuntimeConfig.Group,
		"cluster":     "DEFAULT",
		"healthyOnly": true,
	}, &serviceResponse)
	if err == nil {
		instanceResponse = serviceResponse.Data.Hosts
	}
	return
}

func (this *DiscoveryClient) GetAllServicesInfo() (serviceList response.ServiceList, err error) {
	// 使用web API 获取所有group的服务
	res := response.Result[response.ServiceList]{}
	err = this.getRequestHolder().GET(this.parseUrl("/ns/service/list"), map[string]interface{}{
		"pageNo":         1,
		"pageSize":       512,
		"namespaceId":    this.config.ClientConfig.NamespaceId,
		"groupNameParam": this.config.RuntimeConfig.Group,
	}, &res)
	return res.Data, err
}
