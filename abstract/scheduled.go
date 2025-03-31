package abstract

import (
	_interface "github.com/duanchi/min/v2/interface"
)

type Scheduled struct {
	Bean
	_interface.ScheduledInterface
}

func (this *Scheduled) Run(condition string) {}
