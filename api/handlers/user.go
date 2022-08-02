package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/pkg/user"
)

func FindAll() fiber.Handler {
	return func(c *fiber.Ctx) error {

		userRepo, err := user.NewUserRepository()

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(&fiber.Map{
				"status": false,
				"data":   "",
				"error":  err.Error(),
			})
		}

		result, err := userRepo.FindAll()

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(&fiber.Map{
				"status": false,
				"data":   "",
				"error":  err.Error(),
			})
		}

		return c.JSON(&fiber.Map{
			"status": true,
			"data":   result,
			"error":  nil,
		})

	}
}
