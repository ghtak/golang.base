package fiberfx

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

type ErrorHandler = func(*fiber.Ctx, error) error

func NewDefaultErrorHandler() ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return ctx.Status(fiberErr.Code).JSON(fiber.Map{"message": fiberErr.Message})
		}
		return ctx.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
}
