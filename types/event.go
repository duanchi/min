package types

import "reflect"

type ConditionEvent struct {
	Conditions    map[string]bool
	EventListener []reflect.Value
}

type Event struct {
	EventName     string
	Conditions    []string
	EmitCondition string
}
