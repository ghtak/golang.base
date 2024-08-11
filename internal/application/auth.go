package application

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"go.uber.org/fx"
)

type Auth interface {
	Authentication(ctx context.Context) (context.Context, error)
	SelectAuthentication(ctx context.Context, callMeta interceptors.CallMeta) bool
}

type authImpl struct {
}

func (a *authImpl) Authentication(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
func (a *authImpl) SelectAuthentication(
	ctx context.Context,
	callMeta interceptors.CallMeta) bool {
	return false
}

func NewAuth() Auth {
	return &authImpl{}
}

func NewSelectAuthFunc(auth Auth) func(ctx context.Context, callMeta interceptors.CallMeta) bool {
	return auth.SelectAuthentication
}

func NewAuthFunc(auth Auth) auth.AuthFunc {
	return auth.Authentication
}

var ModuleAuth = fx.Module(
	"ModuleAppAuth",
	fx.Provide(NewAuth, NewSelectAuthFunc, NewAuthFunc),
)
