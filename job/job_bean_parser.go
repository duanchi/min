package job

import (
	"reflect"

	"github.com/duanchi/min/v2/types"
	"github.com/duanchi/min/v2/util"
)

type JobBeanParser struct {
	types.BeanParser
}

func (parser JobBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	if util.IsBeanKind(tag, "job") {
		name := tag.Get("job")

		if name == "" {
			name = beanName
		}

		JobList[name] = bean
	}
}
