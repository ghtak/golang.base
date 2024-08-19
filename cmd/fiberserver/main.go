package main

import (
	"github.com/ghtak/golang.grpc.base/cmd/fiberserver/app/echo"
	"github.com/ghtak/golang.grpc.base/cmd/fiberserver/app/ws"
	"github.com/ghtak/golang.grpc.base/internal/adapter/fiberfx"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewMiddlewares(logger *zap.Logger) fiberfx.Middlewares {
	return fiberfx.MiddlewaresFunc(func(app *fiber.App) error {
		app.Use(fiberzap.New(fiberzap.Config{Logger: logger}))
		app.Use(cors.New(cors.ConfigDefault))
		app.Use(recover.New())
		app.Use("/ws", func(c *fiber.Ctx) error {
			if websocket.IsWebSocketUpgrade(c) {
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		})
		return nil
	})
}

func main() {
	fx.New(
		core.Module,
		fiberfx.Module,
		fx.Provide(
			fiberfx.NewDefaultErrorHandler,
			//fiberfx.NewDefaultMiddlewares,
			NewMiddlewares,
		),
		fx.Provide(
			fiberfx.AsRouter(ws.NewController),
			fiberfx.AsRouter(echo.NewController),
		),
		fx.Invoke(func(result fiberfx.RunResult) {}),
	).Run()
}
