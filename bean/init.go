package bean

import (
	_interface "github.com/duanchi/min/interface"
	"reflect"
)

func Init (rawBean reflect.Value, beanMap map[string]reflect.Value) {
	// AOP
	parseInit(rawBean)
}

func parseInit(rawBean reflect.Value) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	rawBean.Addr().Interface().(_interface.BeanInterface).Init()
}

func parseAop (rawBean reflect.Value) {

}