package grpcadapter

import (
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/fx"
)

var (
	moduleName = "grpcadapter"
)

var Module = fx.Module(
	moduleName,
	fx.Provide(NewServer),
	fx.Provide(core.AsModuleEnv(NewEnv)),
)
