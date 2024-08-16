package gatewayfx

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"strings"
	"time"
)

type GatewayParams struct {
	fx.In
	Env      Env
	Services []Service `group:"gatewayfx.Service"`
}

type GatewayResult struct {
	// fx.Out
	CancelFunc  context.CancelFunc
	Mux         *http.ServeMux
	MuxWithGrpc MuxWithGrpc
}

func NewGateway(p GatewayParams) (GatewayResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	gwMux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	for _, service := range p.Services {
		service.RegisterGateway(ctx, gwMux, p.Env.GrpcAddress, dialOpts)
	}
	// todo
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "gen/openapiv2/helloworld.swagger.json")
	})
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("assets/swagger-ui/dist"))))
	mux.Handle("/", gwMux)
	return GatewayResult{
		Mux:         mux,
		MuxWithGrpc: muxWithGrpc,
		CancelFunc:  cancel,
	}, nil
}

type RunParams struct {
	fx.In
	GatewayResult
	Lc  fx.Lifecycle
	Env Env
}

type RunResult struct {
	Err error
}

func Run(p RunParams) RunResult {
	server := &http.Server{Addr: p.Env.Address, Handler: p.Mux}
	p.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go server.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer p.GatewayResult.CancelFunc()
				defer cancel()
				server.Shutdown(ctx)
				return nil
			},
		})
	return RunResult{
		Err: nil,
	}
}

type MuxWithGrpc func(server *grpc.Server, mux *runtime.ServeMux) http.Handler

var muxWithGrpc = func(s *grpc.Server, gwMux *runtime.ServeMux) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			s.ServeHTTP(w, r)
		} else {
			gwMux.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
