package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/handlers"
	"github.com/mendesbarreto/friday/api/middleware"
)

func RegisterUsersRoutes(app *fiber.App) {
    api := app.Group("/v1/users", middleware.Logger())

	api.Get("/", handlers.UserFindAll())
	app.Post("/", handlers.CreateUser())
	app.Post("/auth", handlers.AuthenticateUser())
}

