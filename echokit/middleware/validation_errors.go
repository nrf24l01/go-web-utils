package middleware

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nrf24l01/go-web-utils/echokit/schemas"

	"github.com/go-playground/validator/v10"
)

// FormatValidationErrors converts validator errors into a slice of schemas.FieldError
func FormatValidationErrors(err error) []schemas.FieldError {
	fieldErrors := make([]schemas.FieldError, 0)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := resolveFieldPath(fieldError)
			tag := fieldError.Tag()

			message := getReadableErrorMessage(fieldName, tag, fieldError.Param())
			fe := schemas.FieldError{
				Field:         fieldName,
				Issue:         message,
				RejectedValue: fieldError.Value(),
			}
			fieldErrors = append(fieldErrors, fe)
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
			fe := schemas.FieldError{
				Field: strings.ToLower(fieldName),
				Issue: message,
			}
			fieldErrors = append(fieldErrors, fe)
		} else {
			fe := schemas.FieldError{
				Field: "error",
				Issue: err.Error(),
			}
			fieldErrors = append(fieldErrors, fe)
		}
	}

	return fieldErrors
}

func resolveFieldPath(fieldError validator.FieldError) string {
	namespace := fieldError.Namespace()
	if namespace == "" {
		return strings.ToLower(fieldError.Field())
	}

	parts := strings.Split(namespace, ".")
	if len(parts) <= 1 {
		return namespace
	}

	structNamespace := fieldError.StructNamespace()
	structParts := strings.Split(structNamespace, ".")
	if len(structParts) == len(parts) {
		parts = parts[1:]
	}

	return strings.Join(parts, ".")
}

// getReadableErrorMessage returns a human-readable error message for a validation tag
func getReadableErrorMessage(fieldName, tag, param string) string {
	messages := map[string]string{
		"required":      fmt.Sprintf("%s is required", fieldName),
		"email":         fmt.Sprintf("%s must be a valid email address", fieldName),
		"min":           fmt.Sprintf("%s must be at least %s characters", fieldName, param),
		"max":           fmt.Sprintf("%s must not exceed %s characters", fieldName, param),
		"numeric":       fmt.Sprintf("%s must be numeric", fieldName),
		"url":           fmt.Sprintf("%s must be a valid URL", fieldName),
		"len":           fmt.Sprintf("%s must be exactly %s characters", fieldName, param),
		"gt":            fmt.Sprintf("%s must be greater than %s", fieldName, param),
		"gte":           fmt.Sprintf("%s must be greater than or equal to %s", fieldName, param),
		"lt":            fmt.Sprintf("%s must be less than %s", fieldName, param),
		"lte":           fmt.Sprintf("%s must be less than or equal to %s", fieldName, param),
		"fromto":        "from must be less than to",
		"maxperiod":     "period must not exceed 90 days",
		"maxfuture":     "timestamp must not be more than 5 minutes in the future",
		"sumrate":       "approvalRate + declineRate must equal 1",
		"groupbyperiod": "period is too large for selected groupBy",
	}

	if message, exists := messages[tag]; exists {
		return message
	}

	return fmt.Sprintf("%s failed validation for %s", fieldName, tag)
}
