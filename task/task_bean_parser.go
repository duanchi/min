package task

import (
	"github.com/duanchi/min/types"
	"reflect"
	"strconv"
)

type TaskBeanParser struct {
	types.BeanParser
}

func (parser TaskBeanParser) Parse (tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {

	isTask, _ := strconv.ParseBool(tag.Get("task"))

	if isTask {
		Tasks = append(Tasks, bean)
	}
}