package main

import (
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcadapter"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func main() {
	fx.New(
		core.Module,
		grpcadapter.Module,
		fx.Invoke(func(*grpc.Server) {}),
	).Run()
}
