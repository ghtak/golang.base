package core

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
	logger *zap.Logger,
) *grpc.Server {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
			// Add any other interceptor you want.
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger(logger), opts...),
			// Add any other interceptor you want.
		),
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

var Module = fx.Module(
	"ModuleCOre",
	fx.Provide(NewEnv, NewNetListener, NewGrpcServer),
	moduleGrpcService,
	moduleZapLogger,
)
