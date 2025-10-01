package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func QueryValidationMiddleware(schemaFactory func() interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			schema := schemaFactory()

			// Bind query parameters
			if err := c.Bind(schema); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid query parameters"})
			}

			// Validate using Echo's validator
			if err := c.Validate(schema); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			}

			// Store validated data in context
			c.Set("validatedQuery", schema)
			return next(c)
		}
	}
}