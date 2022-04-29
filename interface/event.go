package _interface

import "github.com/duanchi/min/types"

type EventInterface interface {
	Conditions() (conditions []string)
	Run(event types.Event, arguments ...interface{})
}
