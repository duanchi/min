package _interface

type RpcInterface interface {
	GetPackageName() (name string)
	SetPackageName(name string)
	GetApplicationName() (name string)
	SetApplicationName(name string)
}
