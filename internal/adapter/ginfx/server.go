package ginfx

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
	"time"
)

type ServerParams struct {
	fx.In
	Lc fx.Lifecycle

	Env         Env
	Middlewares Middlewares `optional:"true"`
}

type ServerResults struct {
	fx.Out

	Engine *gin.Engine
}

func NewServer(p ServerParams) (ServerResults, error) {
	engine := gin.New()
	if p.Middlewares != nil {
		err := p.Middlewares.Use(engine)
		if err != nil {
			return ServerResults{}, err
		}
	}
	srv := &http.Server{Addr: p.Env.ServerAddress(), Handler: engine}
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

	return ServerResults{
		Engine: engine,
	}, nil
}
