package context

import (
	"github.com/duanchi/min/v2/bean"
	"github.com/duanchi/min/v2/config"
	"reflect"
)

type ApplicationContext struct {
}

var applicationContext = new(ApplicationContext)

func GetApplicationContext() *ApplicationContext {
	return applicationContext
}

func (this *ApplicationContext) GetBean(name string) interface{} {
	return bean.Get(name)
}

func (this *ApplicationContext) GetConfig(key string) interface{} {
	return config.Get(key)
}

func (this *ApplicationContext) GetConfigRaw(key string) reflect.Value {
	return config.GetRaw(key)
}
