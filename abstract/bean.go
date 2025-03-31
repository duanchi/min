package abstract

import (
	_interface "github.com/duanchi/min/v2/interface"
)

type Bean struct {
	_interface.BeanInterface
	BeanName string
}

func (this *Bean) Init() {}

func (this *Bean) GetName() (name string) {
	return this.BeanName
}

func (this *Bean) SetName(name string) {
	this.BeanName = name
}
