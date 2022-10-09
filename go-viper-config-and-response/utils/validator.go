package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate = &validator.Validate{}

func V() *validator.Validate {
	return validate
}

func init() {
	validate = validator.New()

	validate.RegisterValidation("tha", ValidateThaAlpha)
	validate.RegisterValidation("eng", ValidateEngAlpha)
}

func ToValidationError(err error) validator.ValidationErrors {
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		return validationErr
	}
	return nil
}

func ValidateEngAlpha(fl validator.FieldLevel) bool {
	engRegex := regexp.MustCompile(`^[a-zA-Z -.,]+$`)
	return engRegex.MatchString(fl.Field().String())
}

func ValidateThaAlpha(fl validator.FieldLevel) bool {
	thaiRegex := regexp.MustCompile("^[\\.\u0E01-\u0E3A\u0E3F-\u0E4D\\s]*$")
	return thaiRegex.MatchString(fl.Field().String())
}
