package core

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewNetListener(env Env) *net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(env.ServerAddress))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return &lis
}

func NewGrpcServer(
	lc fx.Lifecycle,
	listener *net.Listener,
	grpcInterceptorParams GrpcInterceptorParams,
) *grpc.Server {
	interceptorsSize := len(grpcInterceptorParams.GrpcServerInterceptors)
	unaryInterceptors := make([]grpc.UnaryServerInterceptor, interceptorsSize)
	streamInterceptors := make([]grpc.StreamServerInterceptor, interceptorsSize)
	for idx, gsi := range grpcInterceptorParams.GrpcServerInterceptors {
		unaryInterceptors[idx] = gsi.Unary()
		streamInterceptors[idx] = gsi.Stream()
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					err := s.Serve(*listener)
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

var moduleCore = fx.Module(
	"ModuleCore",
	fx.Provide(NewEnv, NewNetListener, NewGrpcServer),
)

var Module = fx.Module(
	"ModuleMain",
	moduleCore,
	moduleGrpcService,
	moduleLog,
)
