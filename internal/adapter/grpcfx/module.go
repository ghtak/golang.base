package grpcfx

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var (
	moduleName = "grpcfx"
)

var Module = fx.Module(
	moduleName,
	fx.Provide(NewServer, NewEnv),
	ModuleService,
)

type RunServerParams struct {
	fx.In
	Server *grpc.Server
}
