package discovery

import (
	"github.com/duanchi/min"
	"github.com/duanchi/min/abstract"
	config2 "github.com/duanchi/min/config"
	"github.com/duanchi/min/event"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/microservice/discovery/nacos/request"
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/types/config"
	"github.com/duanchi/min/util"
	"strconv"
	"time"
)

type instance struct {
	Ip       string
	Port     uint64
	Metadata map[string]string
}

type RegisterHolder struct {
	applicationConfig config.Application
	discoveryConfig   config.Discovery
	discoveryClients  map[string]_interface.DiscoveryInterface
	timeTicker        *time.Ticker
	instance          instance
}

func (this *RegisterHolder) RegisterInstance() {
	registerInstance := request.RegisterInstance{
		Ip:          this.instance.Ip,
		Port:        this.instance.Port,
		Weight:      this.discoveryConfig.Weight,
		Enable:      true,
		Healthy:     true,
		Metadata:    this.instance.Metadata,
		ClusterName: "DEFAULT",
		ServiceName: this.applicationConfig.Name,
		GroupName:   this.discoveryConfig.Group,
		Ephemeral:   true,
	}
	for _, discoveryInstance := range this.discoveryClients {
		discoveryInstance.RegisterInstance(registerInstance)
	}
}

func (this *RegisterHolder) DeregisterInstance() {
	for _, discoveryInstance := range Discovery {
		discoveryInstance.DeregisterInstance(request.DeregisterInstance{
			Ip:          this.instance.Ip,
			Port:        this.instance.Port,
			Cluster:     "DEFAULT",
			ServiceName: this.applicationConfig.Name,
			GroupName:   this.discoveryConfig.Group,
			Ephemeral:   true,
		})
	}
}

func (this *RegisterHolder) StartHeartBeat() {
	interval := config2.Get("Discovery.Client.HeartBeatInterval").(int64)
	this.timeTicker = time.NewTicker(time.Duration(interval) * time.Millisecond)
	for {
		select {
		case <-this.timeTicker.C:
			for _, discoveryInstance := range this.discoveryClients {
				go discoveryInstance.HeartBeat(request.HeartBeat{
					ServiceName: this.applicationConfig.Name,
					GroupName:   this.discoveryConfig.Group,
					Ephemeral:   true,
					Beat: request.BeatInfo{
						Ip:          this.instance.Ip,
						Port:        this.instance.Port,
						Weight:      this.discoveryConfig.Weight,
						ServiceName: this.applicationConfig.Name,
						Cluster:     "DEFAULT",
						Metadata:    this.instance.Metadata,
						Scheduled:   true,
					},
				})
			}
		}
	}
}

func (this *RegisterHolder) StopHeartBeat() {
	this.timeTicker.Stop()
}

func NewRegisterHolder(applicationConfig config.Application, discoveryConfig config.Discovery, discoveryClients map[string]_interface.DiscoveryInterface) (holder *RegisterHolder) {

	port := discoveryConfig.Client.Port
	if port == "" {
		port = applicationConfig.ServerPort

	}
	uintPort, err := strconv.ParseUint(port, 10, 0)
	if err != nil {
		min.Log.Error("Discovery register failed, Invalid port. %s", err.Error())
	}

	ip := discoveryConfig.Client.Ip
	if ip == "" {
		ip = util.GetIp()
	}

	metadata := discoveryConfig.Client.Metadata

	if discoveryConfig.Client.InstanceId != "" {
		metadata["instance-id"] = discoveryConfig.Client.InstanceId
	}

	holder = &RegisterHolder{
		applicationConfig: applicationConfig,
		discoveryConfig:   discoveryConfig,
		discoveryClients:  discoveryClients,
		instance: instance{
			Ip:       ip,
			Port:     uintPort,
			Metadata: metadata,
		},
	}

	holder.RegisterInstance()

	event.On("EXIT", DeregisterInstanceEvent{
		holder: holder,
	}.EventInterface)

	return
}

type DeregisterInstanceEvent struct {
	abstract.Event
	holder *RegisterHolder
}

func (this *DeregisterInstanceEvent) Run(event types.Event, arguments ...interface{}) {
	this.holder.DeregisterInstance()
}
