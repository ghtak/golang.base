package fiberfx

import (
	"context"
	"fmt"
	"github.com/gofiber/contrib/fiberzap/v2"
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
}

type ServerResults struct {
	fx.Out
	App *fiber.App
}

func NewServer(params ServerParams) (ServerResults, error) {
	var app *fiber.App
	if params.ErrorHandler != nil {
		app = fiber.New(fiber.Config{
			ErrorHandler: params.ErrorHandler,
		})
	} else {
		app = fiber.New()
	}
	app.Use(fiberzap.New(fiberzap.Config{Logger: params.Logger}))
	params.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				address := fmt.Sprintf("%s:%d", params.Env.Address, params.Env.Port)
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
