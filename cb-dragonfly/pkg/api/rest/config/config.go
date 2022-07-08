package config

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	pkgconfig "github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/labstack/echo/v4"
	_ "github.com/mitchellh/mapstructure"

	coreconfig "github.com/cloud-barista/cb-dragonfly/pkg/api/core/config"
)

// SetMonConfig 모니터링 정책 설정
// @Summary Set monitoring config
// @Description 모니터링 정책 설정
// @Tags [Setting] Multi-Cloud Monitor Policy Setting
// @Accept  json
// @Produce  json
// @Param monitorInfo body pkgconfig.Monitoring true "Details for an Monitor object"
// @Success 200 {object} pkgconfig.Monitoring
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /config [put]
func SetMonConfig(c echo.Context) error {
	params := pkgconfig.Monitoring{}
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	if (params == pkgconfig.Monitoring{}) {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("Invalid parameter, parameter not defined")))
	}

	monConfig, errCode, err := coreconfig.SetMonConfig(params)
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, monConfig)
}

// GetMonConfig 모니터링 정책 조회
// @Summary Get monitoring config
// @Description 모니터링 정책 조회
// @Tags [Setting] Multi-Cloud Monitor Policy Setting
// @Accept  json
// @Produce  json
// @Success 200 {object} pkgconfig.Monitoring
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /config [get]
func GetMonConfig(c echo.Context) error {
	monConfig, errCode, err := coreconfig.GetMonConfig()
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, monConfig)
}

// ResetMonConfig 모니터링 정책 초기화
// @Summary Reset monitoring config
// @Description 모니터링 정책 초기화
// @Tags [Setting] Multi-Cloud Monitor Policy Setting
// @Accept  json
// @Produce  json
// @Success 200 {object} pkgconfig.Monitoring
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /config/reset [put]
func ResetMonConfig(c echo.Context) error {
	monConfig, errCode, err := coreconfig.ResetMonConfig()
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, monConfig)
}
