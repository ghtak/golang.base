package fiberfx

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	moduleName = "fiberfx"
)

var Module = fx.Module(
	moduleName,
	fx.Provide(NewServer, NewEnv),
	ModuleRouter,
)

type RunServerParams struct {
	fx.In
	App *fiber.App
}
