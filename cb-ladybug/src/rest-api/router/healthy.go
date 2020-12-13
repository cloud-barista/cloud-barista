package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Healthy Method
// @Tags Default
// @Summary Health Check
// @Description for health check
// @ID Healthy
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /healthy [get]
func Healthy(c echo.Context) error {
	return c.String(http.StatusOK, "cb-barista cb-ladybug")
}
