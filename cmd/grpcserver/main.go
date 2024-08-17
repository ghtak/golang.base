package main

import (
	"github.com/ghtak/golang.grpc.base/cmd/grpcserver/app/greeter"
	"github.com/ghtak/golang.grpc.base/cmd/grpcserver/app/user"
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcfx"
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcfx/gatewayfx"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func main() {
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
			//fx.Invoke(func(p grpcfx.RunResult) {}),
		),
		fx.Module(
			"gateway",
			gatewayfx.Module,
			fx.Provide(func(server *grpc.Server) gatewayfx.MultiplexerMiddleware {
				return gatewayfx.NewMultiplexers(
					[]gatewayfx.MultiplexerMiddleware{
						gatewayfx.NewSwaggerMultiplexer(
							"/swagger.json",
							"gen/openapiv2/helloworld.swagger.json",
							"/swagger-ui/",
							"assets/swagger-ui/dist"),
					},
				)
			}),
			fx.Provide(
				gatewayfx.AsService(greeter.NewService),
				gatewayfx.AsService(user.NewService),
			),
			//fx.Invoke(func(r gatewayfx.RunResult) {}),
			fx.Invoke(func(r gatewayfx.RunWithGrpcResult) {}),
		),
	).Run()
}
