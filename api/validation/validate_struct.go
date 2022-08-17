package validation

import (
	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

type ValidationError struct {
	Message string        `json:"message"`
	Errors  []*FieldError `json:"errors"`
}

var validate = validator.New()

func ValidateStruct[Type any](value Type) *ValidationError {
	var errors []*FieldError

	err := validate.Struct(value)

	if err != nil {
		for _, validationError := range err.(validator.ValidationErrors) {
			var element FieldError
			element.Field = validationError.StructField()
			element.Tag = validationError.Tag()
			element.Value = validationError.Param()

			errors = append(errors, &element)
		}

		return &ValidationError{
			Message: err.Error(),
			Errors:  errors,
		}
	}

	return nil

}
