package core

import "go.uber.org/fx"

var Module = fx.Module(
	"core",
	fx.Provide(
		NewEnvLoader,
		fx.Annotate(NewEnv, fx.ParamTags(`group:"ModuleEnv"`))),
)
