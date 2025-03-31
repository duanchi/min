package event

import (
	_interface "github.com/duanchi/min/v2/interface"
	"github.com/duanchi/min/v2/types"
	"reflect"
)

var EventList = make(map[string]types.ConditionEvent)

func On(eventName string, eventHandler _interface.EventInterface) {
	/*EventList[eventName] = types.ConditionEvent{
		Conditions:    map[string]bool{
			"DEFAULT": false,
		},
		EventListener: eventHandleFunction,
	}*/
	AddListener(eventName, eventHandler)
}

func AddListener(eventName string, eventHandler _interface.EventInterface) {
	conditionList := eventHandler.Conditions()
	conditions := map[string]bool{
		"DEFAULT": false,
	}

	if conditionList != nil && len(conditionList) != 0 {
		delete(conditions, "DEFAULT")
		for _, condition := range conditionList {
			conditions[condition] = false
		}
	}

	eventList := []reflect.Value{}
	event := types.ConditionEvent{
		Conditions:    conditions,
		EventListener: []reflect.Value{},
	}
	if _, ok := EventList[eventName]; ok {
		eventList = EventList[eventName].EventListener
	}

	eventList = append(eventList, reflect.ValueOf(eventHandler))
	event.EventListener = eventList
	EventList[eventName] = event
}

func Emit(eventName string, arguments ...interface{}) {
	if eventHandler, ok := EventList[eventName]; ok {
		conditions := make([]string, len(EventList[eventName].Conditions))

		for condition, _ := range EventList[eventName].Conditions {
			conditions = append(conditions, condition)
		}

		event := types.Event{
			EventName:     eventName,
			Conditions:    conditions,
			EmitCondition: "",
		}

		for _, scheduled := range eventHandler.EventListener {
			go scheduled.Interface().(_interface.EventInterface).Emit(event, arguments...)
		}
	}
}

func CommitCondition(eventName string, condition string, arguments ...interface{}) {
	if _, ok := EventList[eventName]; ok {
		if _, has := EventList[eventName].Conditions[condition]; has {
			EventList[eventName].Conditions[condition] = true
			for _, status := range EventList[eventName].Conditions {
				if !status {
					return
				}
			}
			Emit(eventName, arguments...)
		}
	}
}

func RevokeCondition(eventName string, condition string) {
	if _, ok := EventList[eventName]; ok {
		if _, has := EventList[eventName].Conditions[condition]; has {
			EventList[eventName].Conditions[condition] = false
		}
	}
}
