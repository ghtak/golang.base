package grpcadapter

import "go.uber.org/fx"

var Module = fx.Module(
	"grpcadapter",
	fx.Provide(NewServer),
)
