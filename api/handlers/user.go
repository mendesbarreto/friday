package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/dto"
	"github.com/mendesbarreto/friday/api/validation"
	"github.com/mendesbarreto/friday/pkg/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UserFindAll() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		userRepo, err := user.NewUserRepository()

		if err != nil {
			return dto.InternalServerError(err.Error())
		}

		result, err := userRepo.FindAll()

		if err != nil {
			return dto.NotFound(err.Error())
		}

		return ctx.JSON(result)

	}
}

func CreateUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userRepo, err := user.NewUserRepository()

		if err != nil {
			return dto.InternalServerError(err.Error())
		}

		var userRequest dto.CreateUserRequestBody
		err = ctx.BodyParser(&userRequest)

		if err != nil {
			log.Fatal(err)
			return dto.BadRequest(ctx, err.Error())
		}

		validationErr := validation.ValidateStruct(userRequest)

		if validationErr != nil {
			return dto.BadRequestWithValidationError(ctx, validationErr)
		}

		user := user.User{
			Username:  userRequest.Email,
			ID:        primitive.NewObjectID(),
			Password:  userRequest.Password,
			CreatedAt: time.Now(),
		}

		userRepo.Create(&user)

		return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{})
	}
}
