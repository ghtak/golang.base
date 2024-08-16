package ginfx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

const (
	tagRouter string = `group:"ginfx.Router"`
)

type Router interface {
	Register(engine *gin.Engine) error
}

func AsRouter(i interface{}) interface{} {
	return fx.Annotate(i, fx.As(new(Router)), fx.ResultTags(tagRouter))
}

type initRouter struct {
}

func RegisterRouter(routers []Router, engine *gin.Engine) initRouter {
	for _, router := range routers {
		router.Register(engine)
	}
	return initRouter{}
}

var ModuleRouter = fx.Module(
	fmt.Sprintf("%s.Router", moduleName),
	fx.Provide(fx.Annotate(RegisterRouter, fx.ParamTags(tagRouter))),
	fx.Invoke(func(initRouter) {}),
)
