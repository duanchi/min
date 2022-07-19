package context

import (
	"github.com/duanchi/min/bean"
	"github.com/duanchi/min/config"
)

type ApplicationContext struct {
}

func NewApplicationContext() ApplicationContext {
	return ApplicationContext{}
}

func (this *ApplicationContext) GetBean(name string) interface{} {
	return bean.Get(name)
}

func (this *ApplicationContext) GetConfig(key string) interface{} {
	return config.Get(key)
}
