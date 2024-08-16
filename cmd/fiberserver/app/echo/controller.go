package echo

import (
	pb "github.com/ghtak/golang.grpc.base/gen/go"
	"github.com/ghtak/golang.grpc.base/internal/adapter/fiberfx"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	fiberfx.Router
	Echo(ctx *fiber.Ctx) error
	Hello(ctx *fiber.Ctx) error
}

func NewController() Controller {
	return controller{}
}

type controller struct {
}

func (c controller) Register(router fiber.Router) error {
	echo := router.Group("/echo")
	echo.Get("/hello", c.Hello)
	echo.Get("/:echo", c.Echo)
	return nil
}

func (c controller) Echo(ctx *fiber.Ctx) error {
	return ctx.SendString(ctx.Params("echo"))
}

func (c controller) Hello(ctx *fiber.Ctx) error {
	resp := pb.HelloRequest{Name: "xxx"}

	return ctx.JSON(resp)
}
