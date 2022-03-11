package route

import (
	"github.com/duanchi/min/types"
	"reflect"
)

type RouteBeanParser struct {
	types.BeanParser
}

func (parser RouteBeanParser) Parse(tag reflect.StructTag, kind string, bean reflect.Value, definition reflect.Type, beanName string) {

	route := tag.Get("route")
	method := tag.Get("method")

	if route != "" {
		BaseRoutes[route+"@"+method] = bean
	}
}
