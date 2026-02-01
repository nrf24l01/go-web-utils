package schemas

import (
	"github.com/labstack/echo/v4"
)

func GenInternalServerError(c echo.Context) ApiError {
	return GenError(c, INTERNAL_SERVER_ERROR, "Internal Server Error", nil)
}

func GenNotFoundError(c echo.Context) ApiError {
	return GenError(c, NOT_FOUND, "Not Found", nil)
}
