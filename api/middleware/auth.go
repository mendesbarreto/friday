package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mendesbarreto/friday/api/dto"
)

func Protected() fiber.Handler {
	return decodeJwt
}

func cleanToken(bearerText string) string {
	tokenStrings := strings.Split(bearerText, " ")

	if len(tokenStrings) < 1 {
		fmt.Println(tokenStrings[0])
		return tokenStrings[0]
	}

	fmt.Println(tokenStrings[1])
	return tokenStrings[1]
}

func decodeJwt(ctx *fiber.Ctx) error {
	bearerString := ctx.Get("Authorization")

	if len(bearerString) == 0 {
		return dto.BadRequest(ctx, "Missing JWT token")
	}

	tokenString := cleanToken(bearerString)

	_, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("TopSecrete"), nil
	})

	if err == nil {
		return ctx.Next()
	}

	if validationErr, ok := err.(*jwt.ValidationError); ok {
		if validationErr.Errors&jwt.ValidationErrorMalformed != 0 {
			return dto.BadRequest(ctx, "Malformad JWT token")
		} else if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
			return dto.NotAuthorized(ctx, "Token Expired")
		} else if validationErr.Errors&jwt.ValidationErrorNotValidYet != 0 {
			return dto.NotAuthorized(ctx, "Token not valid yet")
		}

		return dto.InternalServerError(ctx, err.Error())
	}

	return dto.InternalServerError(ctx, "Problem to validate jwt")
}
