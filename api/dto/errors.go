package dto

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/validation"
)

func InternalServerError(msg string) error {
	return fiber.NewError(fiber.StatusInternalServerError, msg)
}

func BadRequest(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(fiber.StatusBadRequest).SendString(msg)
}

func BadRequestWithValidationError(ctx *fiber.Ctx, validationError *validation.ValidationError) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(validationError)
}

func NotFound(msg string) error {
	return fiber.NewError(fiber.StatusNotFound, msg)
}
