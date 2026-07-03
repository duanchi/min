package _interface

import (
	"google.golang.org/grpc"
)

type GRPCServerInterface interface {
	GetServiceDesc() *grpc.ServiceDesc
}
