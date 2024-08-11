package core

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type GrpcService interface {
	Register(s *grpc.Server)
}

func AsGrpcService(s interface{}) interface{} {
	return fx.Annotate(
		s,
		fx.As(new(GrpcService)),
		fx.ResultTags(`group:"GrpcService"`))
}

type RegisterServiceInitCtx struct {
}

func RegisterService(services []GrpcService, s *grpc.Server) *RegisterServiceInitCtx {
	for _, svc := range services {
		svc.Register(s)
	}
	return nil
}

var moduleGrpcService = fx.Module(
	"ModuleGrpcService",
	fx.Provide(fx.Annotate(RegisterService, fx.ParamTags(`group:"GrpcService"`))),
	fx.Invoke(func(*RegisterServiceInitCtx) {}),
)
