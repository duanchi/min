package service

import (
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/util"
	"reflect"
)

type FeignBeanParser struct {
	types.BeanParser
}

func (parser FeignBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {

	if util.IsBeanKind(tag, "feign") {
		expression := tag.Get("feign")
		FeignServiceBeans[expression] = bean
	}
}
