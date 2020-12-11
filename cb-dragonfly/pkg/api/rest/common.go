package rest

import (
	"github.com/labstack/echo/v4"
)

func SetMessage(msg string) echo.Map {
	responseMsg := echo.Map{}
	responseMsg["message"] = msg
	return responseMsg
}
