package abstract

import (
	_interface "github.com/duanchi/min/interface"
)

type Scheduled struct {
	Bean
	_interface.TaskInterface
}

func (this *Scheduled) Run() {}
