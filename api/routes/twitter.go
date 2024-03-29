package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/handlers"
	"github.com/mendesbarreto/friday/api/middleware"
)

func RegisterTwitterRoutes(app *fiber.App) {
	api := app.Group("/v1/twitter", middleware.Protected(), middleware.Logger())
	api.Get("/list/:id", handlers.GetTweetsFromToday())
	api.Get("/list/:id/mood", handlers.GetAverageTweetsMoodByList())
}
