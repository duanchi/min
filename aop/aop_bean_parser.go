package aop

import (
	"github.com/duanchi/min/types"
	"reflect"
)

type AopBeanParser struct {
	types.BeanParser
}

func (parser AopBeanParser) Parse (tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {

	rpc := tag.Get("aspect")
	packageName := tag.Get("package")

	if packageName == "" {
		packageName = bean.Elem().Type().PkgPath()
	}

	if rpc != "" {

	}
}
