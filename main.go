package main

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
)

func main() {
    fmt.Print("Hello")

    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, my name is F.R.I.D.A.Y, how can I help you?")
    })



    app.Listen(":3000")
}
