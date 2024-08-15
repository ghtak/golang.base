package grpcfx

import (
	"go.uber.org/fx"
)

var (
	moduleName = "grpcfx"
)

var Module = fx.Module(
	moduleName,
	fx.Provide(NewServer, NewEnv),
	ModuleService,
)
