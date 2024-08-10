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

type RegisterServiceCtx struct {
}

func RegisterService(services []GrpcService, s *grpc.Server) *RegisterServiceCtx {
	for _, svc := range services {
		if t, ok := svc.(interface{ testEmbeddedByValue() }); ok {
			t.testEmbeddedByValue()
		}
		s.RegisterService(svc.ServiceDesc(), svc)
	}
	return nil
}

var moduleGrpcService = fx.Module(
	"ModuleGrpcService",
	fx.Provide(fx.Annotate(RegisterService, fx.ParamTags(`group:"GrpcService"`))),
	fx.Invoke(func(*RegisterServiceCtx) {}))
