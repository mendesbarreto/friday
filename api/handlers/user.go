package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/dto"
	"github.com/mendesbarreto/friday/api/validation"
	"github.com/mendesbarreto/friday/pkg/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

		hash, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)

		user := user.User{
			Username:  userRequest.Email,
			ID:        primitive.NewObjectID(),
			Password:  string(hash),
			CreatedAt: time.Now(),
		}

		userRepo.Create(&user)

		return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{})
	}
}

func AuthenticateUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userRepo, err := user.NewUserRepository()

		if err != nil {
			return dto.InternalServerError(err.Error())
		}

		var authUserBody dto.AuthenticateUserRequestBody
		err = ctx.BodyParser(&authUserBody)

		if err != nil {
			log.Fatal(err)
			return dto.BadRequest(ctx, err.Error())
		}

		user, err := userRepo.FindByUserName(authUserBody.Username)

		if err != nil {
			log.Fatal(err)
			return dto.InternalServerError(err.Error())
		}

		if user == nil {
			userNotFound := dto.NotFound(fmt.Sprintf("User %s not found", authUserBody.Username))
			log.Fatal(userNotFound)
			return userNotFound
		}

		return ctx.SendStatus(fiber.StatusTeapot)

	}
}
