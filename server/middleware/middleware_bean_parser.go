package middleware

import (
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/util"
	"reflect"
)

type MiddlewareBeanParser struct {
	types.BeanParser
}

func (parser MiddlewareBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	if util.IsBeanKind(tag, "middleware") {
		Middlewares = append(Middlewares, bean)
	}
}
