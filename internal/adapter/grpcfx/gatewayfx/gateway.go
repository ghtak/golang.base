package gatewayfx

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func NewGateway(p GatewayParams) (GatewayResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	gwMux := p.SeverMux()
	dialOpts := p.DialOptions()
	for _, service := range p.Services {
		service.RegisterGateway(ctx, gwMux, p.Env.GrpcAddress, dialOpts)
	}
	var handler http.Handler = gwMux
	if p.MultiplexerMiddleware != nil {
		handler = p.MultiplexerMiddleware.Wrap(handler)
	}
	return GatewayResult{
		Handler:    handler,
		CancelFunc: cancel,
	}, nil
}

type GatewayParams struct {
	fx.In
	Env                   Env
	ServerMuxMiddleware   ServerMuxMiddleware   `optional:"true"`
	GrpcDialMiddleware    GrpcDialMiddleware    `optional:"true"`
	MultiplexerMiddleware MultiplexerMiddleware `optional:"true"`
	Services              []Service             `group:"gatewayfx.Service"`
}

type GatewayResult struct {
	// fx.Out
	CancelFunc context.CancelFunc
	Handler    http.Handler
}

func (p *GatewayParams) SeverMux() *runtime.ServeMux {
	if p.ServerMuxMiddleware != nil {
		return runtime.NewServeMux(p.ServerMuxMiddleware.Options()...)
	}
	return runtime.NewServeMux()
}

func (p *GatewayParams) DialOptions() []grpc.DialOption {
	if p.GrpcDialMiddleware != nil {
		return p.GrpcDialMiddleware.Options()
	}
	return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
}
