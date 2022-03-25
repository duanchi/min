package scheduled

import (
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/util"
	"reflect"
)

type ScheduledBeanParser struct {
	types.BeanParser
}

func (parser ScheduledBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	if util.IsBeanKind(tag, "scheduled") {
		expression := tag.Get("scheduled")
		switch expression {
		case "":
		case "@start":
			Scheduled.OnStart = append(Scheduled.OnStart, bean)
		case "@exit":
			Scheduled.OnExit = append(Scheduled.OnExit, bean)
		case "@init":
			Scheduled.OnInit = append(Scheduled.OnInit, bean)
		default:
			Scheduled.Cron = append(Scheduled.Cron, Cron{
				Expression: expression,
				Executor:   bean,
			})
		}
	}
}
