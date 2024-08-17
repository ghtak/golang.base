package fiberfx

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type ServerParams struct {
	fx.In
	ErrorHandler ErrorHandler `optional:"true"`
	Middlewares  Middlewares  `optional:"true"`
}

type ServerResults struct {
	fx.Out
	App    *fiber.App
	Router fiber.Router
}

func NewServer(p ServerParams) (ServerResults, error) {
	var app *fiber.App
	if p.ErrorHandler != nil {
		app = fiber.New(fiber.Config{ErrorHandler: p.ErrorHandler})
	} else {
		app = fiber.New()
	}
	if p.Middlewares != nil {
		p.Middlewares.Use(app)
	}
	return ServerResults{App: app, Router: app}, nil
}
