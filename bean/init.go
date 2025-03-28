package bean

import (
	"fmt"
	_interface "github.com/duanchi/min/interface"
	"reflect"
)

func Init(rawBean reflect.Value, name string, beanMap map[string]reflect.Value) {
	// AOP
	parseInit(rawBean, name)
}

func parseInit(rawBean reflect.Value, name string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[min-framework] Init Bean: Init Error,", rawBean.Type().PkgPath()+"."+name, err)
		} else {
			fmt.Println("[min-framework] Init Bean: " + name)
		}
	}()
	// rawBean.Addr().Interface().(_interface.BeanInterface).SetName(name)
	rawBean.Addr().Interface().(_interface.BeanInterface).Init()
}

func parseAop(rawBean reflect.Value) {

}
