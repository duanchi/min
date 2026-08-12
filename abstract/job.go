package abstract

import _interface "github.com/duanchi/min/v2/interface"

type Job struct {
	Bean
	_interface.JobInterface
}

func (j *Job) Execute(params map[string]any, rawParam string) error {
	return nil
}
