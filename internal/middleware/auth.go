package middleware

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

type AuthMiddleware struct {
}

func (m *AuthMiddleware) Auth(ctx context.Context) (context.Context, error) {
	accessToken, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	ctx = logging.InjectFields(ctx, logging.Fields{"auth", accessToken})
	return context.WithValue(ctx, "access_token", accessToken), nil
}

func (m *AuthMiddleware) Unary() grpc.UnaryServerInterceptor {
	return auth.UnaryServerInterceptor(m.Auth)
}

func (m *AuthMiddleware) Stream() grpc.StreamServerInterceptor {
	return auth.StreamServerInterceptor(m.Auth)
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}
