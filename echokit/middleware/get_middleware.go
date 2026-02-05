package middleware

import (
	"github.com/labstack/echo/v4"
)

func QueryValidationMiddleware(schemaFactory func() interface{}) echo.MiddlewareFunc {
	return ValidationMiddleware(schemaFactory, ValidationOptions{
		Source:     ValidationSourceQuery,
		ContextKey: "validatedQuery",
	})
}
