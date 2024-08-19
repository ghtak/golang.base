package grpcfx

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ServerMiddleware interface {
	Options() []grpc.ServerOption
}

type ServerMiddlewareFunc func() []grpc.ServerOption

func (m ServerMiddlewareFunc) Options() []grpc.ServerOption {
	return m()
}

func NewDefaultServerMiddleware(logger *zap.Logger) ServerMiddleware {
	return ServerMiddlewareFunc(func() []grpc.ServerOption {
		return []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(
				logging.UnaryServerInterceptor(InterceptorLogger(logger), NewLoggingOptions()...),
				//selector.UnaryServerInterceptor(
				//	auth.UnaryServerInterceptor(authFn),
				//	selector.MatchFunc(selectAuthFn)),
				recovery.UnaryServerInterceptor(),
			),
			grpc.ChainStreamInterceptor(
				logging.StreamServerInterceptor(InterceptorLogger(logger), NewLoggingOptions()...),
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

func NewLoggingOptions() []logging.Option {
	return []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall,
			logging.FinishCall,
			logging.PayloadReceived,
			logging.PayloadSent),
	}
}

func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
