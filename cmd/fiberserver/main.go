package main

import (
	"github.com/ghtak/golang.grpc.base/cmd/fiberserver/app/echo"
	"github.com/ghtak/golang.grpc.base/internal/adapter/fiberfx"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		core.Module,
		fiberfx.Module,
		fx.Provide(
			fiberfx.NewDefaultErrorHandler,
			fiberfx.NewDefaultMiddlewares,
		),
		fx.Provide(
			fiberfx.AsRouter(echo.NewController)),
		fx.Invoke(func(app *fiber.App) {}),
	).Run()
}
