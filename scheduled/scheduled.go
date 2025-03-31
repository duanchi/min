package scheduled

import (
	"github.com/duanchi/min/v2/abstract"
	"github.com/duanchi/min/v2/config"
	"github.com/duanchi/min/v2/event"
	_interface "github.com/duanchi/min/v2/interface"
	"github.com/duanchi/min/v2/log"
	"github.com/duanchi/min/v2/types"
	"github.com/robfig/cron/v3"
	"reflect"
)

const (
	RUN_ON_START = "RUN_ON_START"
	RUN_ON_EXIT  = "RUN_ON_EXIT"
	RUN_ON_INIT  = "RUN_ON_INIT"
	RUN_CRON     = "RUN_CRON"
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
		for _, s := range Scheduled.OnStart {
			schedule := s
			schedule.Interface().(_interface.ScheduledInterface).Run(RUN_ON_START)
		}
	}()

	RunCron()
}

func RunOnExit() {
	go func() {
		for _, s := range Scheduled.OnExit {
			schedule := s
			schedule.Interface().(_interface.ScheduledInterface).Run(RUN_ON_EXIT)
		}
	}()
}

func RunOnInit() {
	go func() {
		for _, s := range Scheduled.OnInit {
			schedule := s
			schedule.Interface().(_interface.ScheduledInterface).Run(RUN_ON_INIT)
		}
	}()

}

func RunCron() {
	event.AddListener("EXIT", &CronExitEvent{})
	if len(Scheduled.Cron) > 0 {
		cronInstance = cron.New(cron.WithSeconds())

		for _, s := range Scheduled.Cron {

			schedule := s
			expression := config.Parse(schedule.Expression)

			// expressionValue := reflect.ValueOf(expression)
			// util.ParseValueFromConfigInstance(scheduled.Expression, reflect.Indirect(reflect.ValueOf(expression)), config.ConfigInstance)
			_, err := cronInstance.AddFunc(expression, func() {
				log.Log.Infof("[Scheduled] " + schedule.Executor.Interface().(_interface.BeanInterface).GetName() + " run...")
				schedule.Executor.Interface().(_interface.ScheduledInterface).Run(RUN_CRON)
			})
			if err != nil {
				log.Log.Errorf("[Scheduled] "+schedule.Executor.Interface().(_interface.BeanInterface).GetName()+" init error", err)
			}
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
