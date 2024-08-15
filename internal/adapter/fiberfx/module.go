package fiberfx

import "go.uber.org/fx"

var (
	moduleName = "fiberfx"
)

var Module = fx.Module(
	moduleName,
	fx.Provide(NewServer, NewEnv),
	ModuleRouter,
)
