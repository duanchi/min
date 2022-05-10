package discovery

import (
	"fmt"
	config2 "github.com/duanchi/min/config"
	"github.com/duanchi/min/event"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/microservice/discovery/nacos/request"
	"github.com/duanchi/min/microservice/discovery/nacos/types"
	"github.com/duanchi/min/types/config"
	"github.com/duanchi/min/types/discovery"
	"net/url"
	"strconv"
)

var Discovery map[string]_interface.DiscoveryInterface
var ServiceMap map[string]discovery.Service
var serviceUpdater *ServiceUpdater
var registerHolder *RegisterHolder

func Init() {
	discoveryConfig := config2.Get("Discovery").(config.Discovery)
	applicationConfig := config2.Get("Application").(config.Application)
	discoveryNodes := []*url.URL{}
	discoveryServers := map[string]interface{}{
		"nacos": []types.ServerConfig{},
	}
	for _, nodeDsn := range discoveryConfig.Nodes {
		discoveryUrl, err := url.Parse(nodeDsn)
		if err != nil {
			fmt.Println("[min-framework]: Discovery URL Malformed, \"" + nodeDsn + "\"")
			return
		}
		discoveryNodes = append(discoveryNodes, discoveryUrl)
	}

	for _, discoveryNode := range discoveryNodes {
		switch discoveryNode.Scheme {
		case "nacos":
			scheme := "http"
			if discoveryNode.Query().Get("ssl") == "true" {
				scheme = "https"
			}
			port, _ := strconv.Atoi(discoveryNode.Port())
			discoveryServers["nacos"] = append(discoveryServers["nacos"].([]types.ServerConfig), types.ServerConfig{
				Scheme:      scheme,
				ContextPath: discoveryNode.Path,
				IpAddr:      discoveryNode.Hostname(),
				Port:        uint64(port),
			})
		}
	}

	Discovery = map[string]_interface.DiscoveryInterface{}
	ServiceMap = map[string]discovery.Service{}

	for discoveryType, discoveryServerConfigs := range discoveryServers {
		switch discoveryType {
		case "nacos":
			Discovery["nacos"] = &NacosDiscovery{
				applicationConfig: applicationConfig,
				discoveryConfig:   discoveryConfig,
				clientConfig: types.ClientConfig{
					NamespaceId: discoveryConfig.NamespaceId,
					AppName:     applicationConfig.Name,
					LogLevel:    "debug",
				},
				serverConfig: discoveryServerConfigs.([]types.ServerConfig),
			}
		}

		Discovery[discoveryType].Init()
	}

	if discoveryConfig.Client.Enabled {

	}

	event.Emit("DISCOVERY.INIT")

	fmt.Println("Discovery Update Start!!")

	go func() {
		registerHolder = NewRegisterHolder(applicationConfig, discoveryConfig, Discovery)
		StartHeartBeat()
	}()

	go func() {
		serviceUpdater = NewServiceUpdater(discoveryConfig, Discovery)
		StartServiceUpdater()
	}()
}

func GetServiceList() map[string]discovery.Service {
	return ServiceMap
}

func StartServiceUpdater() {
	serviceUpdater.StartUpdateSchedule()
}

func UpdateService() {
	serviceUpdater.UpdateService()
}

func StopServiceUpdater() {
	serviceUpdater.StopUpdateSchedule()
}

func RegisterInstance(instance request.RegisterInstance) {
	registerHolder.RegisterInstance()
}

func DeregisterInstance(instance request.DeregisterInstance) {
	registerHolder.DeregisterInstance()
}

func StartHeartBeat() {
	registerHolder.StartHeartBeat()
}

func StopHeartBeat() {
	registerHolder.StopHeartBeat()
}
