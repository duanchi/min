package types

import "reflect"

type BaseRoute struct {
	Value  reflect.Value
	Method string
	Path   string
}

type BaseRoutesMap map[string]BaseRoute

type RestfulRoute struct {
	Value       reflect.Value
	ResourceKey string
}

type RestfulRoutesMap map[string]RestfulRoute
