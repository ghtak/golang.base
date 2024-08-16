package gatewayfx

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

const (
	tagService string = `group:"gatewayfx.Service"`
)

type Service interface {
	RegisterGateway(ctx context.Context,
		mux *runtime.ServeMux,
		endpoint string,
		dialOpts []grpc.DialOption)
}

func AsService(s interface{}) interface{} {
	return fx.Annotate(s, fx.As(new(Service)), fx.ResultTags(tagService))
}
