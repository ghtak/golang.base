package main

import (
	"github.com/ghtak/golang.grpc.base/cmd/ginserver/app/echo"
	"github.com/ghtak/golang.grpc.base/internal/adapter/ginfx"
	"github.com/ghtak/golang.grpc.base/internal/core"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		core.Module,
		ginfx.Module,
		fx.Provide(
			ginfx.NewDefaultMiddlewares,
		),
		fx.Provide(
			ginfx.AsRouter(echo.NewController),
		),
		fx.Invoke(func(p ginfx.RunServerParams) {}),
	).Run()
}
