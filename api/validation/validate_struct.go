package validation

import (
	"github.com/go-playground/validator/v10"
)

type ValiationError struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func ValidateStruct[Type any](value Type) []*ValiationError {
	var errors []*ValiationError

	err := validate.Struct(value)

	if err != nil {
		for _, validationError := range err.(validator.ValidationErrors) {
			var element ValiationError
			element.FailedField = validationError.StructField()
			element.Tag = validationError.Tag()
			element.Value = validationError.Param()

			errors = append(errors, &element)
		}
	}
	return errors
}
