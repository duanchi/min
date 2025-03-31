package abstract

import (
	_interface "github.com/duanchi/min/v2/interface"
	"github.com/duanchi/min/v2/types"
)

type Event struct {
	Bean
	_interface.EventInterface
}

func (this *Event) Conditions() (conditions []string) {
	return
}

func (this *Event) Emit(event types.Event, arguments ...interface{}) {}
