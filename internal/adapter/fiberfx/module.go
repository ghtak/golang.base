package fiberfx

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	moduleName = "fiberfx"
)

var Module = fx.Module(
	moduleName,
	fx.Provide(NewServer, NewEnv, RegisterRouter, Run),
)

type RunServerParams struct {
	fx.In
	App *fiber.App
}

type RunParams struct {
	fx.In
	RouterResult
	Lc     fx.Lifecycle
	Env    Env
	Logger *zap.Logger
	App    *fiber.App
}

type RunResult struct {
}

func Run(p RunParams) RunResult {
	p.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go p.App.Listen(p.Env.Address)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return p.App.Shutdown()
			},
		},
	)
	return RunResult{}
}

type DependParams struct {
	fx.In
	ErrorHandler ErrorHandler `optional:"true"`
	Middlewares  Middlewares  `optional:"true"`
}
