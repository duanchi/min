package aop

import "github.com/duanchi/min/abstract"

type Aspect struct {
	Pointcut  Pointcut `execution:"mes.in-mes.io/services/**.Create(..)"`
	abstract.Aspect
}

func Name()  {

}
