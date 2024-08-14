package grpcadapter

import (
	"context"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewServer(lc fx.Lifecycle) *grpc.Server {
	server := grpc.NewServer()

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				lis, err := net.Listen("tcp", "0.0.0.0:9999")
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
				return nil
			},
		})

	return server
}
