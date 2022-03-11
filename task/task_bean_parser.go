package task

import (
	"github.com/duanchi/min/types"
	"reflect"
	"strconv"
)

type TaskBeanParser struct {
	types.BeanParser
}

func (parser TaskBeanParser) Parse(tag reflect.StructTag, kind string, bean reflect.Value, definition reflect.Type, beanName string) {

	isTask := false
	if kind == "middleware" {
		isTask = true
	} else {
		isTask, _ = strconv.ParseBool(tag.Get("middleware"))
	}

	if isTask {
		Tasks = append(Tasks, bean)
	}
}
