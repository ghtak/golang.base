package echo

import (
	"github.com/ghtak/golang.grpc.base/internal/adapter/fiberfx"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	fiberfx.Router
}

func NewController() Controller {
	return controller{}
}

type controller struct {
}

func (c controller) Register(app *fiber.App) error {
	echo := app.Group("/echo")
	echo.Get("/:echo", c.echo)
	return nil
}

func (c controller) echo(ctx *fiber.Ctx) error {
	return ctx.SendString(ctx.Params("echo"))
}
