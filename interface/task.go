package _interface

type TaskInterface interface {
	OnStart()
	OnExit()
	AfterInit()
}