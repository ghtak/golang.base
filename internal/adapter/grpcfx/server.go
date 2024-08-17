package grpcfx

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type ServerParams struct {
	fx.In
	ServerMiddleware ServerMiddleware `optional:"true"`
}

type ServerResults struct {
	fx.Out
	Server *grpc.Server
}

func NewServer(p ServerParams) (ServerResults, error) {
	var server *grpc.Server
	if p.ServerMiddleware != nil {
		server = grpc.NewServer(p.ServerMiddleware.Options()...)
	} else {
		server = grpc.NewServer()
	}
	return ServerResults{Server: server}, nil
}
