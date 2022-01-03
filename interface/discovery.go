package _interface

type DiscoveryInterface interface {
	Init()
	RegisterInstance()
	DeregisterInstance()
	GetService(name string, group string)
	GetServices()
	GetAllInstances()
	GetInstances()
	GetHealthInstance()
	Subscribe()
	UnSubscribe()
}
