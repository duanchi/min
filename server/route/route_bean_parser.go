package route

import (
	serverTypes "github.com/duanchi/min/server/types"
	"github.com/duanchi/min/types"
	"reflect"
)

type RouteBeanParser struct {
	types.BeanParser
}

func (parser RouteBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {

	route := tag.Get("route")
	method := tag.Get("method")

	if route != "" {
		BaseRoutes[route+"@"+method] = serverTypes.BaseRoute{
			Value:  bean,
			Method: method,
			Path:   route,
		}
	}
}
