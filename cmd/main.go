package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/handlers"
	"github.com/mendesbarreto/friday/pkg/infra/database"
	"github.com/mendesbarreto/friday/pkg/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	fmt.Print("Hello")

	app := fiber.New()

	app.Post("/users", func(ctx *fiber.Ctx) error {

		mongo, err := database.New()
		if err != nil {
			log.Fatal(err)
			return err
		}

		user := user.User{
			Username:  "mendesbarreto",
			Id:        primitive.NewObjectID(),
			Password:  "123456",
			CreatedAt: time.Now(),
		}

		usersCollection := mongo.Db.Collection("users")

		usersCollection.InsertOne(ctx.Context(), user)

		return nil
	})

	app.Get("/users", handlers.FindAll())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, my name is F.R.I.D.A.Y, how can I help you?")
	})

	app.Listen(":3000")
}
