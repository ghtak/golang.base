package main

import (
	"github.com/ghtak/golang.grpc.base/cmd/grpcserver/app/greeter"
	"github.com/ghtak/golang.grpc.base/cmd/grpcserver/app/user"
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcfx"
	"github.com/ghtak/golang.grpc.base/internal/adapter/grpcfx/gatewayfx"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServerMiddleware(logger *zap.Logger) grpcfx.ServerMiddleware {
	return grpcfx.ServerMiddlewareFunc(func() []grpc.ServerOption {
		return []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(
				logging.UnaryServerInterceptor(
					grpcfx.InterceptorLogger(logger),
					grpcfx.NewLoggingOptions()...),
				//selector.UnaryServerInterceptor(
				//	auth.UnaryServerInterceptor(authFn),
				//	selector.MatchFunc(selectAuthFn)),
				recovery.UnaryServerInterceptor(),
			),
			grpc.ChainStreamInterceptor(
				logging.StreamServerInterceptor(
					grpcfx.InterceptorLogger(logger),
					grpcfx.NewLoggingOptions()...),
				//selector.StreamServerInterceptor(
				//	auth.StreamServerInterceptor(authFn),
				//	selector.MatchFunc(selectAuthFn)),
				recovery.StreamServerInterceptor(),
			),
			//grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			//	MinTime:             5 * time.Minute,
			//	PermitWithoutStream: false,
			//}),
			//grpc.KeepaliveParams(keepalive.ServerParameters{
			//	MaxConnectionIdle:     15 * time.Minute,
			//	MaxConnectionAge:      30 * time.Minute,
			//	MaxConnectionAgeGrace: 5 * time.Minute,
			//	Time:                  5 * time.Minute,
			//	Timeout:               1 * time.Minute,
			//}),
			//grpc.Creds(...),
		}
	})
}

func main() {
	fx.New(
		core.Module,
		fx.Module(
			"grpc",
			grpcfx.Module,
			fx.Provide(NewServerMiddleware),
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
