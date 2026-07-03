package server

import (
	"reflect"

	"github.com/duanchi/min/v2/log"
	"github.com/duanchi/min/v2/types"
	"github.com/duanchi/min/v2/util"
)

type GrpcServerBeanParser struct {
	types.BeanParser
}

func (parser GrpcServerBeanParser) Parse(tag reflect.StructTag, bean reflect.Value, definition reflect.Type, beanName string) {
	if util.IsBeanKind(tag, "grpc-server") {
		ServerBeans = append(ServerBeans, bean)
		log.Log.Info("[min-framework] Init GRPC Server: " + bean.Elem().Type().Name())
	}
}
