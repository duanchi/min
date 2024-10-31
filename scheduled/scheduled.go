package scheduled

import (
	"fmt"
	"github.com/duanchi/min/abstract"
	"github.com/duanchi/min/config"
	"github.com/duanchi/min/event"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/log"
	"github.com/duanchi/min/types"
	"github.com/robfig/cron/v3"
	"reflect"
)

var Scheduled struct {
	OnStart []reflect.Value
	OnExit  []reflect.Value
	OnInit  []reflect.Value
	Cron    []Cron
}

type Cron struct {
	Expression string
	Executor   reflect.Value
}

var cronInstance *cron.Cron

func Init() {}

func RunOnStart() {
	go func() {
		for _, scheduled := range Scheduled.OnStart {
			scheduled.Interface().(_interface.ScheduledInterface).Run()
		}
	}()
	log.Log.Info("Scheduled has been executed at run on start!")

	RunCron()
}

func RunOnExit() {
	go func() {
		for _, scheduled := range Scheduled.OnExit {
			scheduled.Interface().(_interface.ScheduledInterface).Run()
		}
	}()
	log.Log.Info("Scheduled has been executed at run on exit!")
}

func RunOnInit() {
	go func() {
		for _, scheduled := range Scheduled.OnInit {
			scheduled.Interface().(_interface.ScheduledInterface).Run()
		}
	}()
	log.Log.Info("Scheduled has been executed at run on init!")
}

func RunCron() {
	event.AddListener("EXIT", &CronExitEvent{})
	if len(Scheduled.Cron) > 0 {
		cronInstance = cron.New(cron.WithSeconds())

		for _, scheduled := range Scheduled.Cron {

			expression := config.Parse(scheduled.Expression)
			// expressionValue := reflect.ValueOf(expression)
			// util.ParseValueFromConfigInstance(scheduled.Expression, reflect.Indirect(reflect.ValueOf(expression)), config.ConfigInstance)
			cronInstance.AddFunc(expression, func() {
				scheduled.Executor.Interface().(_interface.ScheduledInterface).Run()
			})
			fmt.Println("Scheduled has been registered!! [" + expression + "]")
		}

		cronInstance.Start()
	}
}

type CronExitEvent struct {
	abstract.Event
}

func (this *CronExitEvent) Run(event types.Event, arguments ...interface{}) {
	cronInstance.Stop()
}
