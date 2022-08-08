package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/handlers"
)

func RegisterUsersRoutes(app *fiber.App) {
	app.Get("/users", handlers.UserFindAll())
	app.Post("/users", handlers.CreateUser())
}
