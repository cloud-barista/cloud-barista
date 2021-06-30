package common

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cloud-barista/cb-tumblebug/src/core/common"
)

// RestInitConfig godoc
// @Summary Init config
// @Description Init config
// @Tags [Admin] System environment
// @Accept  json
// @Produce  json
// @Param configId path string true "Config ID"
// @Success 200 {object} common.ConfigInfo
// @Failure 404 {object} common.SimpleMsg
// @Failure 500 {object} common.SimpleMsg
// @Router /config/{configId} [delete]
func RestInitConfig(c echo.Context) error {
	//id := c.Param("configId")
	if err := Validate(c, []string{"configId"}); err != nil {
		common.CBLog.Error(err)
		return SendMessage(c, http.StatusBadRequest, err.Error())
	}

	err := common.InitConfig(c.Param("configId"))
	if err != nil {
		//mapA := common.SimpleMsg{"Failed to find the config " + id}
		//return c.JSON(http.StatusNotFound, &mapA)
		return SendMessage(c, http.StatusOK, "Failed to init the config "+c.Param("configId"))
	} else {
		//return c.JSON(http.StatusOK, &res)
		return SendMessage(c, http.StatusOK, "The config "+c.Param("configId")+" has been initialized.")
	}
}

// RestGetConfig godoc
// @Summary Get config
// @Description Get config
// @Tags [Admin] System environment
// @Accept  json
// @Produce  json
// @Param configId path string true "Config ID"
// @Success 200 {object} common.ConfigInfo
// @Failure 404 {object} common.SimpleMsg
// @Failure 500 {object} common.SimpleMsg
// @Router /config/{configId} [get]
func RestGetConfig(c echo.Context) error {
	//id := c.Param("configId")
	if err := Validate(c, []string{"configId"}); err != nil {
		common.CBLog.Error(err)
		return SendMessage(c, http.StatusBadRequest, err.Error())
	}

	res, err := common.GetConfig(c.Param("configId"))
	if err != nil {
		//mapA := common.SimpleMsg{"Failed to find the config " + id}
		//return c.JSON(http.StatusNotFound, &mapA)
		return SendMessage(c, http.StatusOK, "Failed to find the config "+c.Param("configId"))
	} else {
		//return c.JSON(http.StatusOK, &res)
		return Send(c, http.StatusOK, res)
	}
}

// Response structure for RestGetAllConfig
type RestGetAllConfigResponse struct {
	//Name string     `json:"name"`
	Config []common.ConfigInfo `json:"config"`
}

// RestGetAllConfig godoc
// @Summary List all configs
// @Description List all configs
// @Tags [Admin] System environment
// @Accept  json
// @Produce  json
// @Success 200 {object} RestGetAllConfigResponse
// @Failure 404 {object} common.SimpleMsg
// @Failure 500 {object} common.SimpleMsg
// @Router /config [get]
func RestGetAllConfig(c echo.Context) error {

	var content RestGetAllConfigResponse

	configList, err := common.ListConfig()
	if err != nil {
		//mapA := common.SimpleMsg{"Failed to list configs."}
		//return c.JSON(http.StatusNotFound, &mapA)
		return SendMessage(c, http.StatusOK, "Failed to list configs.")
	}

	if configList == nil {
		//return c.JSON(http.StatusOK, &content)
		return Send(c, http.StatusOK, content)
	}

	// When err == nil && resourceList != nil
	content.Config = configList
	//return c.JSON(http.StatusOK, &content)
	return Send(c, http.StatusOK, content)

}

// RestPostConfig godoc
// @Summary Create or Update config
// @Description Create or Update config (SPIDER_REST_URL, DRAGONFLY_REST_URL, ...)
// @Tags [Admin] System environment
// @Accept  json
// @Produce  json
// @Param config body common.ConfigReq true "Key and Value for configuration"
// @Success 200 {object} common.ConfigInfo
// @Failure 404 {object} common.SimpleMsg
// @Failure 500 {object} common.SimpleMsg
// @Router /config [post]
func RestPostConfig(c echo.Context) error {

	u := &common.ConfigReq{}
	if err := c.Bind(u); err != nil {
		//return err
		return SendMessage(c, http.StatusBadRequest, err.Error())
	}

	fmt.Println("[Creating or Updating Config]")
	content, err := common.UpdateConfig(u)
	if err != nil {
		//common.CBLog.Error(err)
		////mapA := common.SimpleMsg{"Failed to create the config " + u.Name}
		//mapA := common.SimpleMsg{err.Error()}
		//return c.JSON(http.StatusFailedDependency, &mapA)
		return SendMessage(c, http.StatusBadRequest, err.Error())
	}
	//return c.JSON(http.StatusCreated, content)
	return Send(c, http.StatusOK, content)

}

// RestInitAllConfig godoc
// @Summary Init all configs
// @Description Init all configs
// @Tags [Admin] System environment
// @Accept  json
// @Produce  json
// @Success 200 {object} common.SimpleMsg
// @Failure 404 {object} common.SimpleMsg
// @Router /config [delete]
func RestInitAllConfig(c echo.Context) error {

	err := common.InitAllConfig()
	if err != nil {
		common.CBLog.Error(err)
		return SendMessage(c, http.StatusBadRequest, err.Error())
	}

	return SendMessage(c, http.StatusOK, "All configs has been initialized.")
}
