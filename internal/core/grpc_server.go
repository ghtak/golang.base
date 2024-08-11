package core

import (
	"context"
	"fmt"
	"github.com/ghtak/golang.grpc.base/internal/middleware"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewGrpcServer(
	lc fx.Lifecycle,
	env Env,
	m middleware.Params,
) *grpc.Server {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(m.UnaryMiddlewares()...),
		grpc.ChainStreamInterceptor(m.StreamMiddlewares()...),
	)
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					lis, err := net.Listen("tcp", fmt.Sprintf(env.GrpcAddress))
					if err != nil {
						log.Fatalf("failed to listen: %v", err)
					}
					err = s.Serve(lis)
					if err != nil {
						log.Fatalf("failed to serve: %v", err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		})
	return s
}

var moduleGrpcServer = fx.Module(
	"ModuleGrpcServer",
	fx.Provide(NewGrpcServer),
)
