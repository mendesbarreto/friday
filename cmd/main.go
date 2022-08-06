package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/dto"
	"github.com/mendesbarreto/friday/api/handlers"
	"github.com/mendesbarreto/friday/pkg/infra/database"
	"github.com/mendesbarreto/friday/pkg/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	app := fiber.New()

	app.Post("/users", func(ctx *fiber.Ctx) error {
		mongo, err := database.New()
		if err != nil {
			log.Fatal(err)
			return err
		}

		var userRequest dto.CreateUserRequestBody

		payload := struct {
			Name string `json:"name"`
		}{}

		ctx.BodyParser(&payload)

		log.Println("payload  = ", payload)
		err = ctx.BodyParser(&userRequest)

		if err != nil {
			log.Fatal(err)
			ctx.Status(http.StatusBadRequest)
			return dto.BadRequest(err.Error())
		}

		log.Println(userRequest)
		user := user.User{
			Username:  userRequest.Email,
			ID:        primitive.NewObjectID(),
			Password:  userRequest.Password,
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
