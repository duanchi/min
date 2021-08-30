package _interface

type ServiceInterface interface {
	GetServiceName() (name string)
	SetServiceName(name string)
}
