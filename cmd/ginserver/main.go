package main

import (
	"github.com/ghtak/golang.grpc.base/cmd/ginserver/app/echo"
	"github.com/ghtak/golang.grpc.base/cmd/ginserver/app/ws"
	"github.com/ghtak/golang.grpc.base/internal/adapter/ginfx"
	"github.com/ghtak/golang.grpc.base/internal/core"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

func NewMiddlewares(logger *zap.Logger) ginfx.Middlewares {
	return ginfx.MiddlewaresFunc(func(engine *gin.Engine) error {
		engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		//config := cors.DefaultConfig()
		//config.AllowAllOrigins = true
		//config.AllowCredentials = true
		//engine.Use(cors.New(config))
		engine.Use(ginzap.RecoveryWithZap(logger, true))
		return nil
	})
}

func main() {
	fx.New(
		core.Module,
		ginfx.Module,
		fx.Provide(
			NewMiddlewares,
		),
		fx.Provide(
			ginfx.AsRouter(echo.NewController),
			ginfx.AsRouter(ws.NewController),
		),
		fx.Invoke(func(r ginfx.RunResult) {}),
	).Run()
}
