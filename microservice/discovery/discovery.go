package discovery

import (
	"fmt"
	config2 "github.com/duanchi/min/config"
	"github.com/duanchi/min/event"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/types/config"
	"github.com/duanchi/min/types/discovery"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"net/url"
	"strconv"
)

var Discovery map[string]_interface.DiscoveryInterface

var ServiceMap map[string]discovery.Service

func Init() {
	discoveryConfig := config2.Get("Discovery").(config.Discovery)
	applicationConfig := config2.Get("Application").(config.Application)
	discoveryNodes := []*url.URL{}
	discoveryServers := map[string]interface{}{
		"nacos": []constant.ServerConfig{},
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
			port, _ := strconv.Atoi(discoveryNode.Port())
			discoveryServers["nacos"] = append(discoveryServers["nacos"].([]constant.ServerConfig), *constant.NewServerConfig(discoveryNode.Hostname(), uint64(port), func(config *constant.ServerConfig) {
				if discoveryNode.Query().Get("ssl") == "true" {
					config.Scheme = "https"
				}
			}))
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
				clientConfig: constant.ClientConfig{
					NamespaceId: discoveryConfig.NamespaceId,
					AppName:     applicationConfig.Name,
					LogLevel:    "debug",
				},
				serverConfig: discoveryServerConfigs.([]constant.ServerConfig),
			}
		}

		Discovery[discoveryType].Init()
	}

	event.Emit("DISCOVERY.INIT")
}

func GetServiceList() map[string]discovery.Service {
	return ServiceMap
}
