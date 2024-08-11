package core

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type GrpcServerInterceptor interface {
	Unary() grpc.UnaryServerInterceptor
	Stream() grpc.StreamServerInterceptor
}

func AsGrpcServerInterceptor(i interface{}) interface{} {
	return fx.Annotate(
		i,
		fx.As(new(GrpcServerInterceptor)),
		fx.ResultTags(`group:"GrpcServerInterceptor"`))
}

type GrpcInterceptorParams struct {
	fx.In
	GrpcServerInterceptors []GrpcServerInterceptor `group:"GrpcServerInterceptor"`
}
