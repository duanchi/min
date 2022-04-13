package scheduled

import (
	"fmt"
	"github.com/duanchi/min/config"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/util"
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

func Init() {
}

func RunOnStart() {
	for _, scheduled := range Scheduled.OnStart {
		go scheduled.Interface().(_interface.ScheduledInterface).Run()
	}
	fmt.Println("Scheduled has been executed at run on start!")
}

func RunOnExit() {
	for _, scheduled := range Scheduled.OnExit {
		go scheduled.Interface().(_interface.ScheduledInterface).Run()
	}
	fmt.Println("Scheduled has been executed at run on exit!")
}

func RunOnInit() {
	for _, scheduled := range Scheduled.OnInit {
		go scheduled.Interface().(_interface.ScheduledInterface).Run()
	}
	fmt.Println("Scheduled has been executed at run on init!")
}

func RunCron() {
	if len(Scheduled.Cron) > 0 {
		cronInstance := cron.New(cron.WithSeconds())
		defer cronInstance.Stop()

		for _, scheduled := range Scheduled.Cron {

			expression := ""
			util.ParseValueFromConfigInstance(scheduled.Expression, reflect.ValueOf(expression), config.ConfigInstance)
			cronInstance.AddFunc(expression, scheduled.Executor.Interface().(_interface.ScheduledInterface).Run)
			fmt.Println("Scheduled has been registered!! [" + expression + "]")
		}

		cronInstance.Start()
	}
}
