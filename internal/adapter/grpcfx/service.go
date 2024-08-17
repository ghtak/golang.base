package grpcfx

import (
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

type ServiceParams struct {
	fx.In
	Server   *grpc.Server
	Services []Service `group:"grpcfx.Service"`
}

type ServiceResults struct{}

func RegisterService(p ServiceParams) ServiceResults {
	for _, service := range p.Services {
		service.Register(p.Server)
	}
	return ServiceResults{}
}
