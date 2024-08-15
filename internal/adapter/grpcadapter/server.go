package grpcadapter

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func NewServer(
	env Env,
	log *zap.Logger,
	lc fx.Lifecycle,
) *grpc.Server {
	server := grpc.NewServer()
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				address := fmt.Sprintf("%s:%d", env.Address, env.Port)
				log.Info("grpc start", zap.String("address", address))
				lis, err := net.Listen("tcp", address)
				if err != nil {
					log.Error("failed to listen", zap.Error(err))
					return err
				}
				go func() {
					err = server.Serve(lis)
					if err != nil {
						log.Error("failed to listen", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				server.GracefulStop()
				return nil
			},
		})

	return server
}
