package server

import (
	"net"
	"reflect"
	"strconv"

	_interface "github.com/duanchi/min/v2/interface"
	"github.com/duanchi/min/v2/log"
	"github.com/duanchi/min/v2/types/config"
	"google.golang.org/grpc"
)

var ServerBeans = []reflect.Value{}

func Init(config config.GrpcServer) {
	if len(ServerBeans) == 0 {
		return
	}

	lis, _ := net.Listen("tcp", ":"+strconv.Itoa(config.Port))
	s := grpc.NewServer()

	for _, server := range ServerBeans {
		if t, ok := server.Interface().(interface{ testEmbeddedByValue() }); ok {
			t.testEmbeddedByValue()
		}
		s.RegisterService(server.Interface().(_interface.GRPCServerInterface).GetServiceDesc(), server.Interface())
	}

	go func() {
		err := s.Serve(lis)
		if err != nil {
			log.Log.Errorf("grpc server error: %v", err)
		} else {
			log.Log.Infof("gRPC 服务启动 :%d", config.Port)
		}

	}()
}
