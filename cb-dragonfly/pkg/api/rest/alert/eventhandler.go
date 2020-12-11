package alert

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	"github.com/cloud-barista/cb-dragonfly/pkg/core/alert/eventhandler"
	"github.com/cloud-barista/cb-dragonfly/pkg/core/alert/types"
)

// 알람 이벤트 핸들러 목록 조회
func ListEventHandler(c echo.Context) error {
	eventType := c.QueryParam("eventType")
	eventHandlerList, err := eventhandler.ListEventHandlers(eventType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, eventHandlerList)
}

// 알람 이벤트 핸들러 상세 조회
func GetEventHandler(c echo.Context) error {
	eventType := c.Param("type")
	eventHandlerName := c.Param("name")
	eventHandler, err := eventhandler.GetEventHandler(eventType, eventHandlerName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, eventHandler)
}

// 알람 이벤트 핸들러 생성
func CreateEventHandler(c echo.Context) error {
	eventHandlerReq, err := setEventHandlerReq(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	eventHandler, err := eventhandler.CreateEventHandler(eventHandlerReq.Type, eventHandlerReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, eventHandler)
}

// 알람 이벤트 핸들러 수정
func UpdateEventHandler(c echo.Context) error {
	eventType := c.Param("type")
	eventHandlerName := c.Param("name")
	eventHandlerReq, err := setEventHandlerReq(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	eventHandler, err := eventhandler.UpdateEventHandler(eventType, eventHandlerName, eventHandlerReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, eventHandler)
}

// 알람 이벤트 핸들러 삭제
func DeleteEventHandler(c echo.Context) error {
	eventType := c.Param("type")
	eventHandlerName := c.Param("name")
	err := eventhandler.DeleteEventHandler(eventType, eventHandlerName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, rest.SetMessage(fmt.Sprintf("delete event handler with name %s successfully", eventHandlerName)))
}

func setEventHandlerReq(c echo.Context) (types.AlertEventHandlerReq, error) {
	eventType := c.FormValue("type")
	eventHandlerName := c.FormValue("name")

	// Parameters for Slack
	url := c.FormValue("url")
	channel := c.FormValue("channel")

	// parameters for SMTP
	host := c.FormValue("host")
	from := c.FormValue("from")
	to := c.FormValue("to")
	username := c.FormValue("username")
	password := c.FormValue("password")

	eventHandlerReq := types.AlertEventHandlerReq{
		Name: eventHandlerName,
		Type: eventType,

		// Slack
		Url:     url,
		Channel: channel,

		// SMTP
		Host:     host,
		From:     from,
		Username: username,
		Password: password,
	}

	if eventType == eventhandler.SMTPType {
		// Set port parameters
		port, err := strconv.Atoi(c.FormValue("port"))
		if err != nil {
			return types.AlertEventHandlerReq{}, err
		}
		eventHandlerReq.Port = port
		// Set to array parameters
		if strings.ContainsAny(to, ",") {
			toArr := strings.Split(to, ",")
			eventHandlerReq.To = toArr
		}
	}
	return eventHandlerReq, nil
}
