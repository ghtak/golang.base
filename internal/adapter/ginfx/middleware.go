package ginfx

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type Middlewares interface {
	Use(engine *gin.Engine) error
}

type defaultMiddlewares func(engine *gin.Engine) error

func (h defaultMiddlewares) Use(engine *gin.Engine) error {
	return h(engine)
}

func NewDefaultMiddlewares(logger *zap.Logger) Middlewares {
	return defaultMiddlewares(func(engine *gin.Engine) error {
		engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		engine.Use(ginzap.RecoveryWithZap(logger, true))
		return nil
	})
}
