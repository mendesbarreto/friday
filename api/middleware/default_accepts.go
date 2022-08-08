package middleware

import "github.com/gofiber/fiber/v2"

func DefaultAccepts(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")
	return ctx.Next()
}
