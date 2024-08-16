package grpcfx

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"grpcfx",
	fx.Provide(
		NewServer, NewEnv,
		fx.Annotate(registerService, fx.ParamTags(tagService))),
)

type RunServerParams struct {
	fx.In
	RunRegisterServiceResults
}
