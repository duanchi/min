package rpc

import (
	"fmt"
	"github.com/duanchi/min/types"
	"reflect"
)

type RpcBeanParser struct {
	types.BeanParser
}

func (parser RpcBeanParser) Parse(tag reflect.StructTag, kind string, bean reflect.Value, definition reflect.Type, beanName string) {

	rpc := tag.Get("rpc")
	packageName := tag.Get("package")

	if packageName == "" {
		packageName = bean.Elem().Type().PkgPath()
	}

	if rpc != "" {
		RpcBeans[packageName+"."+bean.Elem().Type().Name()] = struct {
			Package  string
			Instance reflect.Value
		}{Package: packageName, Instance: bean}
		fmt.Println("[min-framework] Init RPC: " + packageName + "." + bean.Elem().Type().Name())
	}
}
