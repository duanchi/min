package _interface

import "github.com/duanchi/min/types/discovery"

type DiscoveryInterface interface {
	Init()
	RegisterInstance()
	DeregisterInstance()
	GetService(name string, group string) (discoveryService discovery.Service, err error)
	GetServiceList() map[string]discovery.Service
	GetAllInstances(serviceName string, group string)
	GetInstances(serviceName string, group string)
	GetHealthInstance()
	Subscribe()
	UnSubscribe()
}
