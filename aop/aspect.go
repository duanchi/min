package aop

type Aspect struct {
	Pointcut  Pointcut `execution:"mes.in-mes.io/services/**.Create(..)"`
}

func Name()  {

}
