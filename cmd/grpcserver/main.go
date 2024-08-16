package main

import (
	"github.com/ghtak/golang.grpc.base/cmd/grpcserver/app"
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcfx"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		core.Module,
		grpcfx.Module,
		fx.Provide(grpcfx.NewDefaultServerMiddlewares),
		fx.Provide(
			grpcfx.AsService(app.NewGreetService),
			grpcfx.AsService(app.NewUserService),
		),
		fx.Invoke(func(params grpcfx.RunServerParams) {}),
	).Run()
}
