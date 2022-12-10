package handlers

import (
	"fmt"

	"github.com/cdipaolo/sentiment"
	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/dto"
	"github.com/mendesbarreto/friday/api/validation"
)

type SentimentResponse struct {
	Score uint8 `json:"score"`
}

type SentimentRequestBody struct {
	Sentence string `json:"sentence"`
}

func GetTextSentiment(text string) (uint8, error) {
	model, err := sentiment.Restore()

	if err != nil {
		return 0, err
	}

	analysis := model.SentimentAnalysis(text, sentiment.English)
	return analysis.Score, nil
}

func GetSentiment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		var body SentimentRequestBody
		err := ctx.BodyParser(&body)

		if err != nil {
			dto.BadRequest(ctx, "The request body has something wrong")
		}

		validationErr := validation.ValidateStruct(&body)

		if validationErr != nil {
			return dto.BadRequestWithValidationError(ctx, validationErr)
		}

		sentiment, err := GetTextSentiment(body.Sentence)

		if err != nil {
			return dto.InternalServerError(ctx, err.Error())
		}


		result := SentimentResponse{
			Score: sentiment,
		}

		fmt.Println(result)

		return ctx.Status(fiber.StatusOK).JSON(result)
	}
}
