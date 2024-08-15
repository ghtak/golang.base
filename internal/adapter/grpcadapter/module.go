package grpcadapter

import (
	"go.uber.org/fx"
)

var (
	moduleName = "grpcadapter"
)

var Module = fx.Module(
	moduleName,
	fx.Provide(NewServer, NewEnv),
)
