package grpcfx

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

var Module = fx.Module(
	"grpcfx",
	fx.Provide(NewServer, NewEnv, Run, RegisterService),
)

type RunParams struct {
	fx.In
	ServiceResults

	Lc     fx.Lifecycle
	Env    Env
	Log    *zap.Logger
	Server *grpc.Server
}

type RunResult struct{}

func Run(p RunParams) (RunResult, error) {
	p.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				p.Log.Info("grpc start", zap.String("address", p.Env.Address))
				lis, err := net.Listen("tcp", p.Env.Address)
				if err != nil {
					p.Log.Error("failed to listen", zap.Error(err))
					return err
				}
				go func() {
					err := p.Server.Serve(lis)
					if err != nil {
						p.Log.Error("gprc serve error", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				p.Server.GracefulStop()
				return nil
			},
		})
	return RunResult{}, nil
}

type DependParams struct {
	fx.In
	ServerMiddleware ServerMiddleware `optional:"true"`
}
