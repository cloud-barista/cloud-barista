package router

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/cloud-barista/cb-mcks/src/core/app"
	"github.com/cloud-barista/cb-mcks/src/core/service"
	"github.com/labstack/echo/v4"
	logger "github.com/sirupsen/logrus"
)

// ListSpec godoc
// @Tags Mcir
// @Summary List Specs
// @Description List Specs
// @ID List Spec
// @Accept json
// @Produce json
// @Param	connection	path	string		true  "Connection Name"
// @Param   control-plane  query    string    	true  "string enums"       Enums(Y, N)
// @Param   cpu-min     query     int        	false  "if Control-Plane, >= 2"           minimum(1)
// @Param   cpu-max     query     int        	false  " <= 99999"          minimum(1)    maximum(99999)
// @Param   memory-min     query     int        false  " if Control-Plane, >= 2"          minimum(1)
// @Param   memory-max     query     int        false  " <= 99999"          minimum(1)    maximum(99999)
// @Success 200 {object} service.SpecList
// @Failure 400 {object} app.Status
// @Router /mcir/connections/{connection}/specs [get]
func ListSpec(c echo.Context) error {

	controlPlane := c.QueryParam("control-plane")
	if controlPlane == "" {
		controlPlane = "N"
	}
	cpumin := validateSpec(c, "cpu-min")
	memorymin := validateSpec(c, "memory-min")
	cpumax := validateSpec(c, "cpu-max")
	memorymax := validateSpec(c, "memory-max")

	lookupSpecs, err := service.VerifySpecList(c.Param("connection"), controlPlane, cpumin, cpumax, memorymin, memorymax)

	if err != nil {
		logger.Warnf("(ListSpec) %s'", err.Error())
		return app.SendMessage(c, http.StatusNotFound, err.Error())
	}

	return app.Send(c, http.StatusOK, lookupSpecs)
}

func validateSpec(c echo.Context, param string) int {
	controlPlane := c.QueryParam("control-plane")
	if c.QueryParam(param) == "" {
		if strings.Contains(param, "min") {
			if controlPlane == "Y" {
				return 2
			} else {
				return 1
			}
		}
		if strings.Contains(param, "max") {
			return 99999
		}
	}
	returnParam, _ := strconv.Atoi(c.QueryParam(param))
	return returnParam
}
