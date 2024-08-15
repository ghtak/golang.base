package grpcadapter

import (
	"context"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewServer(
	env *core.Env,
	lc fx.Lifecycle,
) *grpc.Server {
	moduleEnv := env.ModuleEnvs[moduleName].(Env)
	server := grpc.NewServer()
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				lis, err := net.Listen("tcp", moduleEnv.GrpcadapterAddress)
				if err != nil {
					log.Fatalf("failed to listen: %v", err)
					return err
				}
				go func() {
					err = server.Serve(lis)
					if err != nil {
						log.Fatalf("failed to serve: %v", err)
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
