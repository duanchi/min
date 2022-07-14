package abstract

import (
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/types"
)

type Event struct {
	Bean
	_interface.EventInterface
}

func (this *Event) Conditions() (conditions []string) {
	return
}

func (this *Event) Emit(event types.Event, arguments ...interface{}) {}
