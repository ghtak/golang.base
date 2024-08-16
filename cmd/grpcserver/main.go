package main

import (
	"github.com/ghtak/golang.grpc.base/cmd/grpcserver/app/greeter"
	"github.com/ghtak/golang.grpc.base/cmd/grpcserver/app/user"
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcfx"
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcfx/gatewayfx"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/fx"
)

func main() {
	//fx.New(
	//	core.Module,
	//	grpcfx.Module,
	//	fx.Provide(grpcfx.NewDefaultServerMiddleware),
	//	fx.Provide(
	//		grpcfx.AsService(greeter.NewService),
	//		grpcfx.AsService(user.NewService),
	//	),
	//	fx.Invoke(func(params grpcfx.RunServerParams) {}),
	//).Run()

	fx.New(
		core.Module,
		fx.Module(
			"grpc",
			grpcfx.Module,
			fx.Provide(grpcfx.NewDefaultServerMiddleware),
			fx.Provide(
				grpcfx.AsService(greeter.NewService),
				grpcfx.AsService(user.NewService),
			),
			fx.Invoke(func(params grpcfx.RunServerParams) {}),
		),
		fx.Module(
			"gateway",
			gatewayfx.Module,
			fx.Provide(
				gatewayfx.AsService(greeter.NewService),
				gatewayfx.AsService(user.NewService),
			),
			fx.Invoke(func(r gatewayfx.RunResult) {}),
		),
	).Run()
}
