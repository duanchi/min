package nacos

import (
	config2 "github.com/duanchi/min/config"
	"github.com/duanchi/min/microservice/discovery/nacos/request"
	"github.com/duanchi/min/microservice/discovery/nacos/response"
	"github.com/duanchi/min/requests/http"
	"time"
)

type DiscoveryClient struct {
	config        request.Client
	requestHolder http.Request
}

func NewDiscoveryClient(discoveryConfig request.Client) DiscoveryClient {
	client := DiscoveryClient{
		config: discoveryConfig,
	}

	return client
}

func (this *DiscoveryClient) RegisterInstance(param request.RegisterInstance) (bool, error) {

}

func (this *DiscoveryClient) DeregisterInstance(param request.DeregisterInstance) (bool, error) {

}

func (this *DiscoveryClient) UpdateInstance(param request.UpdateInstance) (bool, error) {

}

func (this *DiscoveryClient) GetService(param request.GetService) (response.Service, error)

func (this *DiscoveryClient) SelectAllInstances(param request.SelectAllInstances) ([]response.Instance, error)

func (this *DiscoveryClient) SelectInstances(param request.SelectInstances) ([]response.Instance, error)

func (this *DiscoveryClient) SelectOneHealthyInstance(param request.SelectOneHealthInstance) (*response.Instance, error)

func (this *DiscoveryClient) Subscribe(param *request.Subscribe) error

func (this *DiscoveryClient) Unsubscribe(param *request.Subscribe) error

func (this *DiscoveryClient) GetAllServicesInfo(param request.GetAllServiceInfo) (response.ServiceList, error)

type ServiceUpdater struct {
	discoveryClient DiscoveryClient
	clientHolder    request.Client
	timeTicker      *time.Ticker
}

func (this *ServiceUpdater) StartUpdateSchedule() {
	for {
		select {
		case <-this.timeTicker.C:
			go this.UpdateService()
		}
	}
}

func (this *ServiceUpdater) StopUpdateSchedule() {
	this.timeTicker.Stop()
}

func (this *ServiceUpdater) UpdateService() {

	interval := config2.Get("Discovery.UpdateInterval").(int64)

	this.timeTicker = time.NewTicker(time.Duration(interval) * time.Millisecond)
	this.discoveryClient.GetAllServicesInfo(request.GetAllServiceInfo{
		NameSpace: this.discoveryClient.config.RuntimeConfig.NamespaceId,
		GroupName: this.discoveryClient.config.RuntimeConfig.Group,
		PageNo:    1,
		PageSize:  512,
	})
}

func NewServiceUpdater(discoveryConfig DiscoveryClient, clientHolder request.Client) ServiceUpdater {
	return ServiceUpdater{
		clientHolder:    clientHolder,
		discoveryClient: discoveryConfig,
	}
}
