package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func TranslateValidationErrors(err validator.ValidationErrors) []string {
	var errorMessages []string

	for _, err := range err {
		var message string

		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("The field '%s' is required.", err.Field())
		case "email":
			message = fmt.Sprintf("The field '%s' must be a valid email address.", err.Field())
		case "nefield":
			message = fmt.Sprintf("The field '%s' must not be the same as field '%s'.", err.Field(), err.Param())
		default:
			message = fmt.Sprintf("Validation error on field '%s' with tag '%s'.", err.Field(), err.Tag())
		}

		errorMessages = append(errorMessages, message)
	}

	return errorMessages
}
