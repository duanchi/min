package _interface

import (
	"github.com/duanchi/min/microservice/discovery/nacos/request"
	"github.com/duanchi/min/microservice/discovery/nacos/response"
	"github.com/duanchi/min/types/discovery"
)

type DiscoveryInterface interface {
	Init()
	RegisterInstance(instance request.RegisterInstance)
	DeregisterInstance(instance request.DeregisterInstance)
	HeartBeat(heartBeat request.HeartBeat)
	GetService(serviceName string) (discoveryService discovery.Service, err error)
	GetServiceList() map[string]discovery.Service
	GetAllInstances(serviceName string) (instances []response.Instance, err error)
	GetInstances(serviceName string) (instances []response.Instance, err error)
}
