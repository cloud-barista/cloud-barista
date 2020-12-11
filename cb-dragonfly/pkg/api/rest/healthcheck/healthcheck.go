package healthcheck

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// API 서버 헬스체크
func Ping(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}

