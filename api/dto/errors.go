package dto

import (
	"github.com/gofiber/fiber/v2"
)

func InternalServerError(msg string) error {
	return fiber.NewError(fiber.StatusInternalServerError, msg)
}

func BadRequest(msg string) error {
	return fiber.NewError(fiber.StatusNotFound, msg)
}

func NotFound(msg string) error {
	return fiber.NewError(fiber.StatusNotFound, msg)
}
