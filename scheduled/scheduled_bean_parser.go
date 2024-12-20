package scheduled

import (
	"fmt"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/types"
	"github.com/duanchi/min/util"
	"reflect"
	"strings"
)

type ScheduledBeanParser struct {
	types.BeanParser
}

func (parser ScheduledBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	if util.IsBeanKind(tag, "scheduled") {
		expressions := strings.Split(tag.Get("scheduled"), ",")

		for _, expression := range expressions {
			expression = strings.TrimSpace(expression)
			switch expression {
			case "":
			case "@start":
				Scheduled.OnStart = append(Scheduled.OnStart, bean)
				fmt.Println("[min-framework] Scheduled " + bean.Interface().(_interface.BeanInterface).GetName() + " has been registry at run on start!")
			case "@exit":
				Scheduled.OnExit = append(Scheduled.OnExit, bean)
				fmt.Println("[min-framework] Scheduled " + bean.Interface().(_interface.BeanInterface).GetName() + " has been registry at run on edit!")
			case "@init":
				Scheduled.OnInit = append(Scheduled.OnInit, bean)
				fmt.Println("[min-framework] Scheduled " + bean.Interface().(_interface.BeanInterface).GetName() + " has been registry at run on init!")
			default:
				Scheduled.Cron = append(Scheduled.Cron, Cron{
					Expression: expression,
					Executor:   bean,
				})
				fmt.Println("[min-framework] Scheduled " + bean.Interface().(_interface.BeanInterface).GetName() + " has been registry at run at " + expression + "!")
			}
		}

	}
}
