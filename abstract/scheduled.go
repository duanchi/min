package abstract

import (
	_interface "github.com/duanchi/min/interface"
)

type Scheduled struct {
	Bean
	_interface.ScheduledInterface
}

func (this *Scheduled) Run(condition string) {}
