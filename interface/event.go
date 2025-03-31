package _interface

import "github.com/duanchi/min/v2/types"

type EventInterface interface {
	Conditions() (conditions []string)
	Emit(event types.Event, arguments ...interface{})
}
