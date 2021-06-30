package healthcheck

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// API 서버 헬스체크
func Ping(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}
