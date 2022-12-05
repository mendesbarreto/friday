package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/cdipaolo/sentiment"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mendesbarreto/friday/api/dto"
	"github.com/mendesbarreto/friday/api/validation"
	userpkg "github.com/mendesbarreto/friday/pkg/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func UserFindAll() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		userRepo, err := userpkg.NewUserRepository()

		if err != nil {
			return dto.InternalServerError(ctx, err.Error())
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
		userRepo, err := userpkg.NewUserRepository()

		if err != nil {
			return dto.InternalServerError(ctx, err.Error())
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

		user, err := userRepo.FindByUserName(userRequest.Email)	
		if err != nil {
			return dto.InternalServerError(ctx, err.Error())
		}

		if user != nil {
			return dto.Conflict(ctx, "User already exists")
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			return dto.InternalServerError(ctx, err.Error())
		}

		newUser := userpkg.User {
			Username:  userRequest.Email,
			ID:        primitive.NewObjectID(),
			Password:  string(hash),
			CreatedAt: time.Now(),
		}

		err = userRepo.Create(&newUser)

		if err != nil {
			return dto.InternalServerError(ctx, err.Error()) 
		}

		return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{})
	}
}

func AuthenticateUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userRepo, err := userpkg.NewUserRepository()

		if err != nil {
			return dto.InternalServerError(ctx, err.Error())
		}

		var authUserBody dto.AuthenticateUserRequestBody
		err = ctx.BodyParser(&authUserBody)

		if err != nil {
			return dto.BadRequest(ctx, err.Error())
		}

		user, err := userRepo.FindByUserName(authUserBody.Username)

		if err != nil {
			return dto.InternalServerError(ctx, err.Error())
		}

		if user == nil {
			userNotFound := dto.NotFound("User not found or password is incorrect")
			return userNotFound
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authUserBody.Password))

		if err != nil {
			return dto.NotFound("User not found or password is incorrect")
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": dto.UserResponseBody {
				ID: user.ID.Hex(),
				Email: user.Username,
			},
		})
		
	    tokenStirng, err := token.SignedString([]byte("TopSecrete"))

		if err != nil { 
			return dto.InternalServerError(ctx, err.Error())
		}

		response := dto.AuthResponseBody {
			// TODO: Add secrete on machine env variables
			Token: tokenStirng,	
		}

		print(token)
		return ctx.JSON(response)

	}
}

type SentimentResponse struct {
   Score uint8 `json:"score"`
}

type SentimentRequestBody struct {
   Sentence string `json:"sentence"`
}

func GetSentiment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		
		model, err := sentiment.Restore()

		if err != nil {
			return dto.InternalServerError(ctx, "The model can not be loaded")
		}

		var body SentimentRequestBody
		err = ctx.BodyParser(&body)

		if err != nil {
			dto.BadRequest(ctx, "The request body has something wrong")
		}

		validationErr := validation.ValidateStruct(&body)

		if validationErr != nil {
			return dto.BadRequestWithValidationError(ctx, validationErr)	
		}

		analysis := model.SentimentAnalysis(body.Sentence, sentiment.English)
			
		result := SentimentResponse{
			Score: analysis.Score,
		}

		fmt.Println(result)

		return ctx.Status(fiber.StatusOK).JSON(result)
	}
}

