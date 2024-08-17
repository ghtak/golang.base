package gatewayfx

import (
	"context"
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcfx"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

var Module = fx.Module(
	"gatewayfx",
	fx.Provide(NewEnv, NewGateway, Run, RunWithGrpc),
)

type RunParams struct {
	fx.In
	GatewayResult
	Lc  fx.Lifecycle
	Env Env
}

type RunResult struct {
	Err error
}

func Run(p RunParams) RunResult {
	server := &http.Server{
		Addr:    p.Env.Address,
		Handler: p.Handler,
	}
	p.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go server.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer p.CancelFunc()
				defer cancel()
				server.Shutdown(ctx)
				return nil
			},
		})
	return RunResult{
		Err: nil,
	}
}

type RunWithGrpcParams struct {
	RunParams
	grpcfx.ServiceResults
	Server *grpc.Server
}

type RunWithGrpcResult struct{}

func RunWithGrpc(p RunWithGrpcParams) RunWithGrpcResult {
	server := &http.Server{
		Addr:    p.Env.Address,
		Handler: NewGrpcMultiplexer(p.Server).Wrap(p.Handler),
	}
	p.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go server.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer p.CancelFunc()
				defer cancel()
				server.Shutdown(ctx)
				return nil
			},
		})
	return RunWithGrpcResult{}
}
