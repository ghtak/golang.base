package fiberfx

import (
	"fmt"
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

type routerInit struct{}

func RegisterRouter(routers []Router, router fiber.Router) routerInit {
	for _, r := range routers {
		r.Register(router)
	}
	return routerInit{}
}

var ModuleRouter = fx.Module(
	fmt.Sprintf("%s.Router", moduleName),
	fx.Provide(fx.Annotate(RegisterRouter, fx.ParamTags(tagRouter))),
	fx.Invoke(func(routerInit) {}),
)
