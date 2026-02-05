package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/go-web-utils/echokit/schemas"
)

func TraceMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uuid, err := uuid.NewV7()
			if err != nil {
				slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ID_GENERATION_FAILED",
					slog.String("error", err.Error()),
				)
				return c.JSON(http.StatusServiceUnavailable, schemas.GenError(c, schemas.INTERNAL_SERVER_ERROR, "REQUEST_ID_GENERATION_FAILED", nil))
			}
			c.Set("traceId", uuid.String())
			c.Set("timestamp", time.Now().UTC().Format(time.RFC3339))
			return next(c)
		}
	}
}
