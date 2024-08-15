package grpcfx

import (
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

const (
	tagService string = `group:"grpcfx.Service"`
)

type Service interface {
	Register(s *grpc.Server)
}

func AsService(s interface{}) interface{} {
	return fx.Annotate(s, fx.As(new(Service)), fx.ResultTags(tagService))
}

type ServiceInit struct{}

func RegisterService(services []Service, s *grpc.Server) ServiceInit {
	for _, service := range services {
		service.Register(s)
	}
	return ServiceInit{}
}

var ModuleService = fx.Module(
	fmt.Sprintf("%s.Service", moduleName),
	fx.Provide(fx.Annotate(RegisterService, fx.ParamTags(tagService))),
	fx.Invoke(func(ServiceInit) {}),
)
