package config

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"

	coreconfig "github.com/cloud-barista/cb-dragonfly/pkg/core/config"
)

// 모니터링 정책 설정
func SetMonConfig(c echo.Context) error {
	params, err := c.FormParams()
	if len(params) == 0 {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("Invalid parameter, parameter not defined")))
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}

	paramsMap := map[string]interface{}{}
	for k, _ := range params {
		v := params.Get(k)
		paramsMap[k], err = strconv.Atoi(v)
		if err != nil || paramsMap[k] == 0 {
			return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("Invalid parameter values, %s=%s", k, v)))
		}
	}

	var newMonConfig config.Monitoring
	err = mapstructure.Decode(paramsMap, &newMonConfig)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	monConfig, errCode, err := coreconfig.SetMonConfig(newMonConfig)
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, monConfig)
}

// 모니터링 정책 조회
func GetMonConfig(c echo.Context) error {
	monConfig, errCode, err := coreconfig.GetMonConfig()
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, monConfig)
}

// 모니터링 정책 초기화
func ResetMonConfig(c echo.Context) error {
	monConfig, errCode, err := coreconfig.ResetMonConfig()
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, monConfig)
}
