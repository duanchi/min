package discovery

import (
	"fmt"
	config2 "github.com/duanchi/min/config"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/types/config"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"net/url"
	"strconv"
)

var Discovery _interface.DiscoveryInterface

func Init() {
	discoveryConfig := config2.Get("Discovery").(config.Discovery)
	applicationConfig := config2.Get("Application").(config.Application)
	discoveryUrl, err := url.Parse(discoveryConfig.Url)
	if err != nil {
		fmt.Println("[min-framework]: Discovery URL Malformed, \"" + discoveryConfig.Url + "\"")
		return
	}

	switch discoveryUrl.Scheme {
	case "nacos":
		scheme := "http"
		port, _ := strconv.Atoi(discoveryUrl.Port())
		if discoveryUrl.Query().Get("ssl") == "true" {
			scheme = "https"
		}

		Discovery = &NacosDiscovery{
			applicationConfig: applicationConfig,
			discoveryConfig:   discoveryConfig,
			clientConfig: constant.ClientConfig{
				NamespaceId: discoveryConfig.NamespaceId,
				AppName:     applicationConfig.Name,
				CacheDir:    "/tmp/cache",
				LogDir:      "/tmp/log",
				MaxAge:      3,
				LogLevel:    "debug",
			},
			serverConfig: constant.ServerConfig{
				Scheme:      scheme,
				ContextPath: discoveryUrl.Path,
				IpAddr:      discoveryUrl.Hostname(),
				Port:        uint64(port),
			},
		}

		Discovery.Init()
	}
}
