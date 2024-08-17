package fiberfx

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

const (
	tagRouter string = `group:"fiberfx.Router"`
)

type Router interface {
	Register(router fiber.Router) error
}

func AsRouter(i interface{}) interface{} {
	return fx.Annotate(i, fx.As(new(Router)), fx.ResultTags(tagRouter))
}

type RouterParams struct {
	fx.In
	Router  fiber.Router
	Routers []Router `group:"fiberfx.Router"`
}

type RouterResult struct {
}

func RegisterRouter(p RouterParams) RouterResult {
	for _, r := range p.Routers {
		r.Register(p.Router)
	}
	return RouterResult{}
}
