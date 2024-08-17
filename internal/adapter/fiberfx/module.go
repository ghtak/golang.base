package fiberfx

import (
	"context"
	"fmt"
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
				address := fmt.Sprintf("%s:%d", p.Env.Address, p.Env.Port)
				go p.App.Listen(address)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return p.App.Shutdown()
			},
		},
	)
	return RunResult{}
}
