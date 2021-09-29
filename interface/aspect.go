package _interface

type AspectInterface interface {
	Before()
	After()
	Around()
	AfterReturning()
	AfterPanic()
}