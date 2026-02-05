package middleware

import (
	"net/http"

	"github.com/nrf24l01/go-web-utils/auth"
	"github.com/nrf24l01/go-web-utils/config"
	"github.com/nrf24l01/go-web-utils/echokit/schemas"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(config config.JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Извлекаем токен из заголовка Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, schemas.GenError(c, schemas.UNAUTHORIZED, "missing token", nil))
			}

			// Убираем "Bearer " из заголовка
			if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
				return c.JSON(http.StatusUnauthorized, schemas.GenError(c, schemas.UNAUTHORIZED, "invalid token format", nil))
			}
			tokenString := authHeader[7:]

			if len(config.AccessJWTSecret) == 0 {
				return c.JSON(http.StatusInternalServerError, schemas.GenInternalServerError(c))
			}

			// Проверяем токен
			claims, err := auth.ValidateToken(tokenString, []byte(config.AccessJWTSecret))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, schemas.GenError(c, schemas.UNAUTHORIZED, "invalid or expired token", nil))
			}

			// Извлекаем user_id
			userID, ok := claims["user_id"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, schemas.GenError(c, schemas.UNAUTHORIZED, "invalid token claims", nil))
			}

			// Передаем user_id в контекст
			c.Set("userID", userID)

			return next(c)
		}
	}
}
