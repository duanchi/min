package _interface

import "github.com/duanchi/min/types/discovery"

type DiscoveryInterface interface {
	Init()
	RegisterInstance()
	DeregisterInstance()
	GetService(name string, group string) (discoveryService discovery.Service, err error)
	GetServices()
	GetAllInstances()
	GetInstances()
	GetHealthInstance()
	Subscribe()
	UnSubscribe()
}
