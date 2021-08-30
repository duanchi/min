package task

import (
	_interface "github.com/duanchi/min/interface"
	"reflect"
)

var Tasks []reflect.Value

func Init () {

}

func RunOnStart() {
	for _, task :=range Tasks {
		go task.Interface().(_interface.TaskInterface).OnStart()
	}
}

func RunOnExit() {
	for _, task :=range Tasks {
		go task.Interface().(_interface.TaskInterface).OnExit()
	}
}

func RunAfterInit() {
	for _, task :=range Tasks {
		go task.Interface().(_interface.TaskInterface).AfterInit()
	}
}
