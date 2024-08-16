package fiberfx

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ServerParams struct {
	fx.In
	Lc           fx.Lifecycle
	Env          Env
	Logger       *zap.Logger
	ErrorHandler ErrorHandler `optional:"true"`
	Middlewares  Middlewares  `optional:"true"`
}

type ServerResults struct {
	fx.Out
	App *fiber.App
}

func NewServer(p ServerParams) (ServerResults, error) {
	var app *fiber.App
	if p.ErrorHandler != nil {
		app = fiber.New(fiber.Config{
			ErrorHandler: p.ErrorHandler,
		})
	} else {
		app = fiber.New()
	}
	if p.Middlewares != nil {
		p.Middlewares.Use(app)
	}
	p.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				address := fmt.Sprintf("%s:%d", p.Env.Address, p.Env.Port)
				go app.Listen(address)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return app.Shutdown()
			},
		},
	)
	return ServerResults{App: app}, nil
}
