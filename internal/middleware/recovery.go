package middleware

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
)

type RecoveryMiddleware struct {
}

func (m RecoveryMiddleware) Unary() grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor()
}

func (m RecoveryMiddleware) Stream() grpc.StreamServerInterceptor {
	return recovery.StreamServerInterceptor()
}

func NewRecoveryMiddleware() RecoveryMiddleware {
	return RecoveryMiddleware{}
}
