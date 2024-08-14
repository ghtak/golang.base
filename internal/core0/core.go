package core0

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"ModuleCore",
	fx.Provide(NewEnv),
	moduleGrpcServer,
	moduleGrpcService,
	moduleLog,
)
