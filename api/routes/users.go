package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/handlers"
)

func RegisterUsersRoutes(app *fiber.App) {
	app.Get("/users", handlers.UserFindAll())
	app.Post("/users", handlers.CreateUser())
	app.Post("/users/auth", handlers.AuthenticateUser())
}

func RegisterMLRoutes(app *fiber.App) {
	app.Post("/text/sentiment", handlers.GetSentiment())
}

func RegisterTwitterRoutes(app *fiber.App) {
	app.Get("/twitter/list/:id", handlers.GetTweetsFromToday())
}
