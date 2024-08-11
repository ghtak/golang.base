package middleware

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewLoggingOptions() []logging.Option {
	return []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall,
			logging.FinishCall,
			logging.PayloadReceived,
			logging.PayloadSent),
	}
}

func NewUnaryServerInterceptors(
	logger *zap.Logger,
	opts []logging.Option,
	authFn auth.AuthFunc,
	selectAuthFn func(ctx context.Context, callMeta interceptors.CallMeta) bool,
) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
		selector.UnaryServerInterceptor(
			auth.UnaryServerInterceptor(authFn),
			selector.MatchFunc(selectAuthFn)),
		recovery.UnaryServerInterceptor(),
	}
}

func NewStreamServerInterceptors(
	logger *zap.Logger,
	opts []logging.Option,
	authFn auth.AuthFunc,
	selectAuthFn func(ctx context.Context, callMeta interceptors.CallMeta) bool,
) []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		logging.StreamServerInterceptor(InterceptorLogger(logger), opts...),
		selector.StreamServerInterceptor(
			auth.StreamServerInterceptor(authFn),
			selector.MatchFunc(selectAuthFn)),
		recovery.StreamServerInterceptor(),
	}
}

var Module = fx.Module(
	"ModuleMiddleware",
	fx.Provide(
		NewLoggingOptions,
		NewUnaryServerInterceptors,
		NewStreamServerInterceptors,
	),
)

type Params struct {
	fx.In
	UnaryServerInterceptors  []grpc.UnaryServerInterceptor
	StreamServerInterceptors []grpc.StreamServerInterceptor
}
