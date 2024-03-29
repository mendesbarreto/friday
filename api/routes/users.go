package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/handlers"
	"github.com/mendesbarreto/friday/api/middleware"
)

func RegisterUsersRoutes(app *fiber.App) {
	api := app.Group("/v1/users", middleware.Logger())
	api.Get("/", middleware.Protected(), handlers.UserFindAll())
	api.Post("/", handlers.CreateUser())
	api.Post("/auth", handlers.AuthenticateUser())
}
