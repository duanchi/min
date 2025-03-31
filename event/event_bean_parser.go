package event

import (
	_interface "github.com/duanchi/min/v2/interface"
	"github.com/duanchi/min/v2/types"
	"github.com/duanchi/min/v2/util"
	"reflect"
)

type EventBeanParser struct {
	types.BeanParser
}

func (parser EventBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	if util.IsBeanKind(tag, "event") {

		conditionList := bean.Interface().(_interface.EventInterface).Conditions()
		conditions := map[string]bool{
			"DEFAULT": false,
		}

		if conditionList != nil && len(conditionList) != 0 {
			delete(conditions, "DEFAULT")
			for _, condition := range conditionList {
				conditions[condition] = false
			}
		}

		eventName := tag.Get("event")
		eventList := []reflect.Value{}
		event := types.ConditionEvent{
			Conditions:    conditions,
			EventListener: []reflect.Value{},
		}
		if _, ok := EventList[eventName]; ok {
			eventList = EventList[eventName].EventListener
		}

		eventList = append(eventList, bean)
		event.EventListener = eventList
		EventList[eventName] = event
	}
}
