package middleware

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Middleware interface {
	Unary() grpc.UnaryServerInterceptor
	Stream() grpc.StreamServerInterceptor
}

func AsMiddleware(i interface{}) interface{} {
	return fx.Annotate(i, fx.As(new(Middleware)), fx.ResultTags(`group:"Middleware"`))
}

type Params struct {
	fx.In
	Middlewares []Middleware `group:"Middleware"`
}

func (p Params) UnaryMiddlewares() []grpc.UnaryServerInterceptor {
	interceptorsSize := len(p.Middlewares)
	interceptors := make([]grpc.UnaryServerInterceptor, interceptorsSize)
	for idx, i := range p.Middlewares {
		interceptors[idx] = i.Unary()
	}
	return interceptors
}

func (p Params) StreamMiddlewares() []grpc.StreamServerInterceptor {
	interceptorsSize := len(p.Middlewares)
	interceptors := make([]grpc.StreamServerInterceptor, interceptorsSize)
	for idx, i := range p.Middlewares {
		interceptors[idx] = i.Stream()
	}
	return interceptors
}

var Module = fx.Module(
	"ModuleMiddleware",
	fx.Provide(
		AsMiddleware(NewLoggingMiddleware),
		AsMiddleware(NewRecoveryMiddleware),
		//AsMiddleware(NewAuthMiddleware),
	),
)
