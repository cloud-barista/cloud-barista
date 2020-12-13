package app

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/validation"
	"github.com/cloud-barista/cb-ladybug/src/core/model"
	"github.com/labstack/echo/v4"
)

type Status struct {
	Message string `json:"message"`
}

func SendMessage(c echo.Context, httpCode int, msg string) error {
	return c.JSON(httpCode, Status{Message: msg})
}

func Send(c echo.Context, httpCode int, json interface{}) error {
	return c.JSON(httpCode, json)
}

func Validate(c echo.Context, params []string) error {
	valid := validation.Validation{}

	for _, name := range params {
		valid.Required(c.Param(name), name)
	}

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(fmt.Sprintf("[%s]%s", err.Key, err.Error()))
		}
	}
	return nil
}

func ClusterReqValidate(c echo.Context, req model.ClusterReq) error {
	if req.ControlPlaneNodeCount != 1 {
		return errors.New("control plane node count must be one")
	}
	if req.WorkerNodeCount < 1 {
		return errors.New("worker node count must be at least one")
	}

	return nil
}

func NodeReqValidate(c echo.Context, req model.NodeReq) error {
	if req.Config == "" {
		return errors.New("config is required")
	}
	if req.WorkerNodeSpec == "" {
		return errors.New("worker node spec is required")
	}
	if req.WorkerNodeCount < 1 {
		return errors.New("worker node count must be at least one")
	}

	return nil
}
