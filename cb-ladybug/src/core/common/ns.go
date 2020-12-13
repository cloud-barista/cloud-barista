package common

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NsValidate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ns := c.Param("namespace")
			if ns == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Namespace cannot be empty")
			}
			return next(c)
		}
	}
}
