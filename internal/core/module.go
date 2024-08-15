package core

import "go.uber.org/fx"

var (
	moduleName = "core"
)

var Module = fx.Module(
	moduleName,
	fx.Provide(
		fx.Annotate(
			NewEnvRepository,
			fx.ParamTags(`group:"NamedEnv"`),
		),
	),
	fx.Provide(AsNamedEnv(NewEnv)),
)
