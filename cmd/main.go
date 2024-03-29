package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/middleware"
	"github.com/mendesbarreto/friday/api/routes"
)

func main() {
	app := fiber.New(fiber.Config{
		DisableKeepalive:  true,
		EnablePrintRoutes: true,
	})

	app.Use(middleware.DefaultAccepts)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, my name is F.R.I.D.A.Y, how can I help you?")
	})

	routes.RegisterUsersRoutes(app)
	routes.RegisterMLRoutes(app)
	routes.RegisterTwitterRoutes(app)

	app.Listen(":3000")
}
