package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/handlers"
	"github.com/mendesbarreto/friday/api/middleware"
)

func RegisterMLRoutes(app *fiber.App) {
    api := app.Group("/v1/ml", middleware.Logger())


    api.Post("/text/sentiment", handlers.GetSentiment())
}

