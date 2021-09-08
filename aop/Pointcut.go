package aop

type Pointcut struct {
	Execution string
	Within string
	Args []string
	This *interface{}
	target *interface{}
}