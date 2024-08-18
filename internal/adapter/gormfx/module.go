package gormfx

import "go.uber.org/fx"

var Module = fx.Module(
	"gormfx",
	fx.Provide(NewDatabase, NewEnv),
)
