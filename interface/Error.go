package _interface

type Error interface {
	error

	Code() int
}
