package task

import (
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/util"
	"reflect"
)

type TaskBeanParser struct {
	types.BeanParser
}

func (parser TaskBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	if util.IsBeanKind(tag, "task") {
		Tasks = append(Tasks, bean)
	}
}
