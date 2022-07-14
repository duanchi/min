package types

import "reflect"

type RestfulRoute struct {
	Value       reflect.Value
	ResourceKey string
}

type RestfulRoutesMap map[string]RestfulRoute
