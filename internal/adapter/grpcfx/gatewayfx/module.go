package gatewayfx

import "go.uber.org/fx"

var Module = fx.Module(
	"gatewayfx",
	fx.Provide(NewEnv, NewGateway, Run),
)
