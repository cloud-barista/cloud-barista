package healthcheck

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Ping API 서버 헬스체크
// @Summary Server Health Check
// @Description 서버 헬스체크
// @Tags [Health] Health Check
// @Accept  json
// @Produce  json
// @Router /healthcheck [get]
func Ping(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}
