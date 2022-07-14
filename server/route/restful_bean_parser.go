package route

import (
	types2 "github.com/duanchi/min/server/types"
	"github.com/duanchi/min/types"
	"reflect"
	"strings"
)

type RestfulBeanParser struct {
	types.BeanParser
}

func (parser RestfulBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {

	resource := tag.Get("restful")

	if resource == "" {
		resource = tag.Get("rest")
	}

	if resource != "" {
		key := tag.Get("resource_key")
		if key == "" {
			key = "id"
		}
		resources := strings.Split(resource, ",")
		for _, res := range resources {
			RestfulRoutes[strings.TrimSpace(res)] = types2.RestfulRoute{
				Value:       bean,
				ResourceKey: key,
			}
		}
	}
}
