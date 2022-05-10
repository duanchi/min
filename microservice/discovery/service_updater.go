package discovery

import (
	"fmt"
	config2 "github.com/duanchi/min/config"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/types/config"
	"time"
)

type ServiceUpdater struct {
	discoveryConfig  config.Discovery
	discoveryClients map[string]_interface.DiscoveryInterface
	timeTicker       *time.Ticker
}

func (this *ServiceUpdater) StartUpdateSchedule() {
	interval := config2.Get("Discovery.UpdateInterval").(int64)
	this.timeTicker = time.NewTicker(time.Duration(interval) * time.Millisecond)
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
	for _, discoveryClient := range this.discoveryClients {
		fmt.Println("RUN Service Updater!!!")
		discoveryClient.GetServiceList()
	}
}

func NewServiceUpdater(discoveryConfig config.Discovery, discoveryClients map[string]_interface.DiscoveryInterface) *ServiceUpdater {
	return &ServiceUpdater{
		discoveryConfig:  discoveryConfig,
		discoveryClients: discoveryClients,
	}
}
