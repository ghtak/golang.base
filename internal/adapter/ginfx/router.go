package ginfx

import (
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

type RouterParams struct {
	fx.In
	Routers []Router `group:"ginfx.Router"`
	Engine  *gin.Engine
}

type RouterResult struct {
}

func RegisterRouter(p RouterParams) RouterResult {
	for _, router := range p.Routers {
		router.Register(p.Engine)
	}
	return RouterResult{}
}
