package cache

import (
	"github.com/duanchi/min/v2/types"
	"reflect"
)

type MiddlewareBeanParser struct {
	types.BeanParser
}

func (parser MiddlewareBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {

	cacheName := tag.Get("cache")

	if cacheName != "" {
		CacheEngines[cacheName] = bean
	}
}
