package _interface

type Error interface {
	Error() string
	Code() int
	Status() int
	Data() interface{}
}
