package gatewayfx

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"strings"
)

type ServerMuxMiddleware interface {
	Options() []runtime.ServeMuxOption
}

type GrpcDialMiddleware interface {
	Options() []grpc.DialOption
}

type MultiplexerMiddleware interface {
	Wrap(handler http.Handler) http.Handler
}

func NewDefaultServerMuxMiddleware() ServerMuxMiddleware {
	return serverMuxMiddlewareFunc(func() []runtime.ServeMuxOption {
		return []runtime.ServeMuxOption{}
	})
}

func NewDefaultSGrpcDialMiddleware() GrpcDialMiddleware {
	return grpcDialMiddlewareFunc(
		func() []grpc.DialOption {
			return []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials())}
		},
	)
}

type grpcDialMiddlewareFunc func() []grpc.DialOption

func (f grpcDialMiddlewareFunc) Options() []grpc.DialOption {
	return f()
}

type serverMuxMiddlewareFunc func() []runtime.ServeMuxOption

func (f serverMuxMiddlewareFunc) Options() []runtime.ServeMuxOption {
	return f()
}

type SwaggerMultiplexer struct {
	jsonUrl  string
	jsonFile string
	uiUrl    string
	uiDir    string
}

func (s SwaggerMultiplexer) Wrap(handler http.Handler) http.Handler {
	sMux := http.NewServeMux()
	sMux.HandleFunc(s.jsonUrl, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, s.jsonFile)
	})
	sMux.Handle(s.uiUrl, http.StripPrefix(s.uiUrl, http.FileServer(http.Dir(s.uiDir))))
	sMux.Handle("/", handler)
	return sMux
}

type GrpcMultiplexer struct {
	server *grpc.Server
}

func NewSwaggerMultiplexer(jsonUrl, jsonFile, uiUrl, uiDir string) MultiplexerMiddleware {
	return SwaggerMultiplexer{jsonUrl: jsonUrl, jsonFile: jsonFile, uiUrl: uiUrl, uiDir: uiDir}
}

func (s GrpcMultiplexer) Wrap(handler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 &&
			strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			s.server.ServeHTTP(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func NewGrpcMultiplexer(server *grpc.Server) MultiplexerMiddleware {
	return GrpcMultiplexer{server: server}
}

type CombinedMultiplexer struct {
	multiplexers []MultiplexerMiddleware
}

func (m CombinedMultiplexer) Wrap(handler http.Handler) http.Handler {
	for _, mux := range m.multiplexers {
		handler = mux.Wrap(handler)
	}
	return handler
}

func NewMultiplexers(multiplexers []MultiplexerMiddleware) MultiplexerMiddleware {
	return CombinedMultiplexer{multiplexers: multiplexers}
}
