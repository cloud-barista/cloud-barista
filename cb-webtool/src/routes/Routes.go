package routes

import (
	"github.com/labstack/echo"
)

type Route struct {
	Method, Path string
	Function     echo.HandlerFunc
}

const (
	CbDriverURL  = "http://localhost"
	CbDriverPort = "1234"
)
