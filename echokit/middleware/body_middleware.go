package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func BodyValidationMiddleware(schemaFactory func() interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			schema := schemaFactory()

			if err := c.Bind(schema); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Invalid request payload"})
			}

			if err := c.Validate(schema); err != nil {
				return c.JSON(http.StatusBadRequest, FormatValidationErrors(err))
			}

			c.Set("validatedBody", schema)
			return next(c)
		}
	}
}
