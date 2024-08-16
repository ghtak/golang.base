package ginfx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var (
	moduleName = "ginfx"
)

var Module = fx.Module(
	moduleName,
	fx.Provide(NewServer, NewEnv),
	ModuleRouter,
)

type RunServerParams struct {
	fx.In
	Engine *gin.Engine
}
