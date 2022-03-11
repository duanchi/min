package service

import (
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/types"
	"reflect"
)

type ServiceBeanParser struct {
	types.BeanParser
}

func (parser ServiceBeanParser) Parse(tag reflect.StructTag, kind string, bean reflect.Value, definition reflect.Type, beanName string) {
	if definition.Implements(reflect.TypeOf((*_interface.ServiceInterface)(nil)).Elem()) {
		bean.Elem().Interface().(_interface.ServiceInterface).SetServiceName(beanName)
	}
}
