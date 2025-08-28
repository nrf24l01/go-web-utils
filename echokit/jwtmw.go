package echokit

import (
	"net/http"

	"github.com/nrf24l01/go-web-utils/jwtutil"

	"github.com/labstack/echo/v4"
)

// JWTMiddleware создает middleware для проверки JWT токена
func JWTMiddleware(accessSecret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Извлекаем токен из заголовка Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
			}

			// Убираем "Bearer " из заголовка
			if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token format"})
			}
			tokenString := authHeader[7:]

			if len(accessSecret) == 0 {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "jwt secret is not configured"})
			}

			// Проверяем токен
			claims, err := jwtutil.ValidateToken(tokenString, accessSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			}

			// Извлекаем user_id
			userID, ok := claims["user_id"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
			}

			// Передаем user_id в контекст
			c.Set("userID", userID)

			return next(c)
		}
	}
}
