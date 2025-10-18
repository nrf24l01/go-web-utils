package middleware

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func BodyValidationMiddleware(schemaFactory func() interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			schema := schemaFactory()

			if err := c.Bind(schema); err != nil {
				// If the bind error comes from validation, return formatted validation errors
				if _, ok := err.(validator.ValidationErrors); ok {
					return c.JSON(http.StatusUnprocessableEntity, FormatValidationErrors(err))
				}

				// For other bind errors return 400 Bad Request
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
			}

			if err := c.Validate(schema); err != nil {
				// Return 422 for validation errors with formatted message
				return c.JSON(http.StatusUnprocessableEntity, FormatValidationErrors(err))
			}

			c.Set("validatedBody", schema)
			return next(c)
		}
	}
}
