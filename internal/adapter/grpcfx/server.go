package grpcfx

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type ServerParams struct {
	fx.In
	Lc                fx.Lifecycle
	Env               Env
	Log               *zap.Logger
	ServerMiddlewares ServerMiddlewares `optional:"true"`
}

type ServerResults struct {
	fx.Out
	Server *grpc.Server
}

func NewServer(p ServerParams) (ServerResults, error) {
	var server *grpc.Server
	if p.ServerMiddlewares != nil {
		server = grpc.NewServer(p.ServerMiddlewares.Options()...)
	} else {
		server = grpc.NewServer()
	}
	p.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				address := fmt.Sprintf("%s:%d", p.Env.Address, p.Env.Port)
				p.Log.Info("grpc start", zap.String("address", address))
				lis, err := net.Listen("tcp", address)
				if err != nil {
					p.Log.Error("failed to listen", zap.Error(err))
					return err
				}
				go server.Serve(lis)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				server.GracefulStop()
				return nil
			},
		})

	return ServerResults{Server: server}, nil
}
