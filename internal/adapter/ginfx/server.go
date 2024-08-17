package ginfx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ServerParams struct {
	fx.In
	Lc fx.Lifecycle

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
	return ServerResults{Engine: engine}, nil
}
