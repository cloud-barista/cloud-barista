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
// @Success 200 {object} rest.SimpleMsg
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /healthcheck [get]
func Ping(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}
