package middleware

import (
	"net/http"
	"time"

	"github.com/nrf24l01/go-web-utils/echokit/schemas"

	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/go-web-utils/config"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

func TGMiddleware(config config.TgWebAppConfig) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{Message: "missing token", Code: http.StatusUnauthorized})
			}

			// Убираем "Bearer " из заголовка
			tokenString := ""
			if len(authHeader) > 4 && authHeader[:4] == "tma " {
				tokenString = authHeader[4:]
			} else {
				return c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{Message: "invalid token format", Code: http.StatusUnauthorized})
			}

			expInHour := time.Duration(config.InitDataExpireHours) * time.Hour

			verifyErr := initdata.Validate(tokenString, config.TgBotToken, expInHour)
			if verifyErr != nil {
				return c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{Message: "invalid token", Code: http.StatusUnauthorized})
			}
			tokenData, err := initdata.Parse(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{Message: "invalid token", Code: http.StatusUnauthorized})
			}
			if tokenData.User.ID != 0 {
				c.Set("userID", tokenData.User.ID)
				c.Set("userName", tokenData.User.Username)
			} else {
				return c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{Message: "token does not contain user data", Code: http.StatusUnauthorized})
			}
			return next(c)
		}
	}
}