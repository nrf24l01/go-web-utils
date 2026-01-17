package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/go-web-utils/echokit/schemas"
)

type ValidationSource int

const (
	ValidationSourceBody ValidationSource = iota
	ValidationSourceQuery
	ValidationSourceMultipart
)

type ValidationOptions struct {
	Source     ValidationSource
	ContextKey string
}

func ValidationMiddleware(schemaFactory func() interface{}, opts ValidationOptions) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			schema := schemaFactory()

			switch opts.Source {
			case ValidationSourceMultipart:
				form, err := c.MultipartForm()
				if err != nil {
					status, payload := buildValidationResponse(c, http.StatusBadRequest, schemas.BAD_REQUEST, "Invalid multipart form", err)
					return c.JSON(status, payload)
				}
				if err := bindMultipartForm(c, schema, form); err != nil {
					status, payload := buildValidationResponse(c, http.StatusBadRequest, schemas.BAD_REQUEST, "Invalid multipart form", err)
					return c.JSON(status, payload)
				}
			case ValidationSourceQuery:
				if err := c.Bind(schema); err != nil {
					status, payload := buildValidationResponse(c, http.StatusBadRequest, schemas.BAD_REQUEST, "Invalid query parameters", err)
					return c.JSON(status, payload)
				}
			default:
				if err := c.Bind(schema); err != nil {
					status, payload := buildValidationResponse(c, http.StatusBadRequest, schemas.BAD_REQUEST, "Invalid request payload", err)
					return c.JSON(status, payload)
				}
			}

			if err := c.Validate(schema); err != nil {
				status, payload := buildValidationResponse(c, http.StatusUnprocessableEntity, schemas.VALIDATION_FAILED, "Validation failed", err)
				return c.JSON(status, payload)
			}

			contextKey := opts.ContextKey
			if contextKey == "" {
				contextKey = "validated"
			}
			c.Set(contextKey, schema)
			return next(c)
		}
	}
}

func buildValidationResponse(c echo.Context, status int, code schemas.ErrorCode, message string, err error) (int, schemas.ValidationError) {
	if err == nil {
		err = errors.New(message)
	}
	apiErr := schemas.GenError(c, code, message, nil)
	return status, schemas.ValidationError{
		ApiError:    apiErr,
		FieldErrors: FormatValidationErrors(err),
	}
}