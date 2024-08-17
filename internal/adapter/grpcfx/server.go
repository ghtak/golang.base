package grpcfx

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type ServerParams struct {
	fx.In
	ServerMiddleware ServerMiddleware `optional:"true"`
}

func NewServer(p ServerParams) *grpc.Server {
	if p.ServerMiddleware != nil {
		return grpc.NewServer(p.ServerMiddleware.Options()...)
	}
	return grpc.NewServer()
}
