package middleware

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FormatValidationErrors converts validator errors into readable error messages
func FormatValidationErrors(err error) map[string]interface{} {
	errors := make(map[string]interface{})

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			tag := fieldError.Tag()

			message := getReadableErrorMessage(fieldName, tag, fieldError.Param())
			errors[strings.ToLower(fieldName)] = message
		}
	} else {
		// Try to parse common validator error string formats, e.g.:
		// "Key: 'Struct.Field' Error:Field validation for 'Field' failed on the 'required' tag"
		re := regexp.MustCompile(`Field validation for '([^']+)' failed on the '([^']+)' tag`)
		matches := re.FindStringSubmatch(err.Error())
		if len(matches) == 3 {
			fieldName := matches[1]
			tag := matches[2]
			message := getReadableErrorMessage(fieldName, tag, "")
			errors[strings.ToLower(fieldName)] = message
		} else {
			errors["error"] = err.Error()
		}
	}

	return map[string]interface{}{
		"errors": errors,
	}
}

// getReadableErrorMessage returns a human-readable error message for a validation tag
func getReadableErrorMessage(fieldName, tag, param string) string {
	messages := map[string]string{
		"required": fmt.Sprintf("%s is required", fieldName),
		"email":    fmt.Sprintf("%s must be a valid email address", fieldName),
		"min":      fmt.Sprintf("%s must be at least %s characters", fieldName, param),
		"max":      fmt.Sprintf("%s must not exceed %s characters", fieldName, param),
		"numeric":  fmt.Sprintf("%s must be numeric", fieldName),
		"url":      fmt.Sprintf("%s must be a valid URL", fieldName),
		"len":      fmt.Sprintf("%s must be exactly %s characters", fieldName, param),
		"gt":       fmt.Sprintf("%s must be greater than %s", fieldName, param),
		"gte":      fmt.Sprintf("%s must be greater than or equal to %s", fieldName, param),
		"lt":       fmt.Sprintf("%s must be less than %s", fieldName, param),
		"lte":      fmt.Sprintf("%s must be less than or equal to %s", fieldName, param),
	}

	if message, exists := messages[tag]; exists {
		return message
	}

	return fmt.Sprintf("%s failed validation for %s", fieldName, tag)
}
