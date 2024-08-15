package fiberfx

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type Middlewares interface {
	Use(app *fiber.App) error
}

type defaultMiddlewares func(app *fiber.App) error

func (h defaultMiddlewares) Use(app *fiber.App) error {
	return h(app)
}

func NewDefaultMiddlewares(logger *zap.Logger) Middlewares {
	return defaultMiddlewares(func(app *fiber.App) error {
		app.Use(fiberzap.New(fiberzap.Config{Logger: logger}))
		app.Use(cors.New(cors.ConfigDefault))
		app.Use(recover.New())
		return nil
	})
}
