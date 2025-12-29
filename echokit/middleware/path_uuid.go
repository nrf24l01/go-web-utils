package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func PathUuidV4Middleware(param string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uuid := c.Param(param)
			if !isValidUuidV4(uuid) {
				return echo.NewHTTPError(400, "Invalid UUIDv4 in path parameter")
			}
			return next(c)
		}
	}
}

func isValidUuidV4(s string) bool {
	u, err := uuid.Parse(s)
	if err != nil {
		return false
	}
	return u.Version() == uuid.Version(4)
}
