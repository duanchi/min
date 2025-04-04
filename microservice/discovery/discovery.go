package discovery

import (
	config2 "github.com/duanchi/min/v2/config"
	"github.com/duanchi/min/v2/event"
	_interface "github.com/duanchi/min/v2/interface"
	"github.com/duanchi/min/v2/log"
	"github.com/duanchi/min/v2/microservice/discovery/nacos/request"
	"github.com/duanchi/min/v2/microservice/discovery/nacos/types"
	"github.com/duanchi/min/v2/types/config"
	"github.com/duanchi/min/v2/types/discovery"
	"net/url"
	"strconv"
	"strings"
)

var Discovery map[string]_interface.DiscoveryInterface
var ServiceMap map[string]discovery.Service
var serviceUpdater *ServiceUpdater
var registerHolder *RegisterHolder

func Init() {
	discoveryConfig := config2.Get("Discovery").(config.Discovery)
	applicationConfig := config2.Get("Application").(config.Application)
	httpServerConfig := config2.Get("HttpServer").(config.HttpServer)
	discoveryNodes := []*url.URL{}
	discoveryServers := map[string][]types.ServerConfig{
		"nacos": {},
	}
	for _, nodeDsn := range strings.Split(discoveryConfig.Nodes, ",") {
		discoveryUrl, err := url.Parse(strings.Trim(nodeDsn, " "))
		if err != nil {
			log.Log.Error("[min-framework]: Discovery URL Malformed, \"" + nodeDsn + "\"")
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
			discoveryServers["nacos"] = append(discoveryServers["nacos"], types.ServerConfig{
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
				serverConfig: discoveryServerConfigs,
			}
		}
		Discovery[discoveryType].Init()
	}

	if discoveryConfig.Client.Enabled {

	}

	event.Emit("DISCOVERY.INIT")

	log.Log.Debugf("Discovery Update Start...")

	go func() {
		registerHolder = NewRegisterHolder(applicationConfig, httpServerConfig, discoveryConfig, Discovery)
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
