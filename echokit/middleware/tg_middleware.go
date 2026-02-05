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
				return c.JSON(http.StatusUnauthorized, schemas.GenError(c, schemas.UNAUTHORIZED, "missing token", nil))
			}

			// Убираем "Bearer " из заголовка
			tokenString := ""
			if len(authHeader) > 4 && authHeader[:4] == "tma " {
				tokenString = authHeader[4:]
			} else {
				return c.JSON(http.StatusUnauthorized, schemas.GenError(c, schemas.UNAUTHORIZED, "invalid token format", nil))
			}

			expInHour := time.Duration(config.InitDataExpireHours) * time.Hour

			verifyErr := initdata.Validate(tokenString, config.TgBotToken, expInHour)
			if verifyErr != nil {
				return c.JSON(http.StatusUnauthorized, schemas.GenError(c, schemas.UNAUTHORIZED, "invalid token", nil))
			}
			tokenData, err := initdata.Parse(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, schemas.GenError(c, schemas.UNAUTHORIZED, "invalid token", nil))
			}
			if tokenData.User.ID != 0 {
				c.Set("userID", tokenData.User.ID)
				c.Set("userName", tokenData.User.Username)
			} else {
				return c.JSON(http.StatusUnauthorized, schemas.GenError(c, schemas.UNAUTHORIZED, "token does not contain user data", nil))
			}
			return next(c)
		}
	}
}
