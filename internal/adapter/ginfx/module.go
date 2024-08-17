package ginfx

import "C"
import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
	"time"
)

var (
	moduleName = "ginfx"
)

var Module = fx.Module(
	moduleName,
	fx.Provide(NewServer, NewEnv, RegisterRouter, Run),
)

type RunParams struct {
	fx.In
	RouterResult
	Lc     fx.Lifecycle
	Env    Env
	Engine *gin.Engine
}

type RunResult struct{}

func Run(p RunParams) RunResult {
	srv := &http.Server{Addr: p.Env.Address, Handler: p.Engine}
	p.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go srv.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				srv.Shutdown(ctx)
				return nil
			},
		})
	return RunResult{}
}
