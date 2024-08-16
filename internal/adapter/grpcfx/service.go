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

type RunRegisterServiceResults struct{}

func registerService(services []Service, s *grpc.Server) RunRegisterServiceResults {
	for _, service := range services {
		service.Register(s)
	}
	return RunRegisterServiceResults{}
}
