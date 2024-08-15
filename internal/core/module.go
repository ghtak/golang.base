package core

import "go.uber.org/fx"

var (
	moduleName = "core"
)

var Module = fx.Module(
	moduleName,
	fx.Provide(NewEnv, NewLogger),
)
