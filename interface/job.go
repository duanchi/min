package _interface

type JobInterface interface {
	Execute(params map[string]any, rawParam string) error
}
