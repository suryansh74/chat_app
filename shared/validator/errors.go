package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func TranslateValidationErrors(validationErr error) []ValidationError {
	if validationErr == nil {
		return nil
	}

	var ve validator.ValidationErrors
	if !errors.As(validationErr, &ve) {
		return nil
	}

	var result []ValidationError
	for _, e := range ve {
		field := cleanFieldName(e.Field())
		message := translateTag(e.Tag(), e.Param(), field)
		result = append(result, ValidationError{
			Field:   field,
			Message: message,
		})
	}

	return result
}

func cleanFieldName(field string) string {
	parts := strings.Split(field, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return field
}

func translateTag(tag, param, field string) string {
	lowerField := strings.ToLower(field)

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", lowerField)
	case "email":
		return fmt.Sprintf("%s must be a valid email", lowerField)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", lowerField, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", lowerField, param)
	case "eqfield":
		return fmt.Sprintf("%s must match %s", lowerField, strings.ToLower(param))
	default:
		return fmt.Sprintf("invalid %s", lowerField)
	}
}
