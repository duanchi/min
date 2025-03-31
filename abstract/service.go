package abstract

import _interface "github.com/duanchi/min/v2/interface"

type Service struct {
	Bean
	_interface.ServiceInterface
}

//func (this *Service) Init() {
//	this.Bean.Init()
//}

func (this *Service) GetServiceName() (name string) {
	return this.BeanName
}

func (this *Service) SetServiceName(name string) {
	this.BeanName = name
}
