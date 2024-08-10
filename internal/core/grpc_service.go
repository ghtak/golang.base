package core

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type GrpcService interface {
	ServiceDesc() *grpc.ServiceDesc
}

func AsGrpcService(s interface{}) interface{} {
	return fx.Annotate(
		s,
		fx.As(new(GrpcService)),
		fx.ResultTags(`group:"GrpcService"`))
}
