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

func NewGrpcServer(lc fx.Lifecycle, listener *net.Listener) *grpc.Server {
	s := grpc.NewServer()
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

var Module = fx.Module(
	"ModuleCOre",
	fx.Provide(NewEnv, NewNetListener, NewGrpcServer),
	moduleGrpcService,
)
