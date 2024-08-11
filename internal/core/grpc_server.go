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

func serveGRPC(s *grpc.Server, env Env) {
	lis, err := net.Listen("tcp", fmt.Sprintf(env.GrpcAddress))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func NewGrpcServer(
	lc fx.Lifecycle,
	env Env,
	middleware middleware.Params,
) *grpc.Server {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.UnaryServerInterceptors...),
		grpc.ChainStreamInterceptor(middleware.StreamServerInterceptors...))
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go serveGRPC(s, env)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				s.GracefulStop()
				return nil
			},
		})
	return s
}

var moduleGrpcServer = fx.Module(
	"ModuleGrpcServer",
	fx.Provide(NewGrpcServer),
)
