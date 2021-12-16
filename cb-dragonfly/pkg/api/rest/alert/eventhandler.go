package alert

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/eventhandler"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	"github.com/labstack/echo/v4"
)

// ListEventHandler 알람 이벤트 핸들러 목록 조회
// @Summary List monitoring alert event-handler
// @Description 알람 이벤트 핸들러 목록 조회
// @Tags [EventHandler] Alarm Event Handler management
// @Accept  json
// @Produce  json
// @Param eventType query string false "이벤트 핸들러 유형" Enums(slack, smtp)
// @Success 200 {object} []types.AlertEventHandler
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /alert/eventhandlers [get]
func ListEventHandler(c echo.Context) error {
	eventType := c.QueryParam("eventType")
	eventHandlerList, err := eventhandler.ListEventHandlers(eventType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, eventHandlerList)
}

// GetEventHandler 알람 이벤트 핸들러 상세 조회
// @Summary Get monitoring alert event-handler
// @Description 알람 이벤트 핸들러 조회
// @Tags [EventHandler] Alarm Event Handler management
// @Accept  json
// @Produce  json
// @Param type path string true "이벤트 핸들러 유형"
// @Param name path string true "이벤트 핸들러 이름"
// @Success 200 {object} types.AlertEventHandler
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /alert/eventhandler/type/{type}/event/{name} [get]
func GetEventHandler(c echo.Context) error {
	eventType := c.Param("type")
	eventHandlerName := c.Param("name")
	eventHandler, err := eventhandler.GetEventHandler(eventType, eventHandlerName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, eventHandler)
}

// CreateEventHandler 알람 이벤트 핸들러 생성
// @Summary Create monitoring alert event-handler
// @Description 알람 이벤트 핸들러 생성
// @Tags [EventHandler] Alarm Event Handler management
// @Accept  json
// @Produce  json
// @Param eventHandlerInfo body types.AlertEventHandlerReq true "Details for an EventHandler object"
// @Success 200 {object} types.AlertEventHandler
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /alert/eventhandler [post]
func CreateEventHandler(c echo.Context) error {
	params := &types.AlertEventHandlerReq{}
	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}

	eventHandler, err := eventhandler.CreateEventHandler(params.Type, *params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, eventHandler)
}

// UpdateEventHandler 알람 이벤트 핸들러 수정
// @Summary Update monitoring alert event-handler
// @Description 알람 이벤트 핸들러 수정
// @Tags [EventHandler] Alarm Event Handler management
// @Accept  json
// @Produce  json
// @Param type path string true "이벤트 핸들러 유형"
// @Param name path string true "이벤트 핸들러 이름"
// @Param eventHandlerInfo body types.AlertEventHandlerReq true "Details for an EventHandler (slack) object"
// @Success 200 {object} types.AlertEventHandler
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /alert/eventhandler/type/{type}/event/{name} [put]
func UpdateEventHandler(c echo.Context) error {
	eventType := c.Param("type")
	eventHandlerName := c.Param("name")
	params := &types.AlertEventHandlerReq{}
	params.Type = eventType
	params.Name = eventHandlerName
	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	eventHandler, err := eventhandler.UpdateEventHandler(eventType, eventHandlerName, *params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, eventHandler)
}

// DeleteEventHandler 알람 이벤트 핸들러 삭제
// @Summary Delete monitoring alert event-handler
// @Description 알람 이벤트 핸들러 삭제
// @Tags [EventHandler] Alarm Event Handler management
// @Accept  json
// @Produce  json
// @Param type path string true "이벤트 핸들러 유형"
// @Param name path string true "이벤트 핸들러 이름"
// @Success 200 {object} rest.SimpleMsg
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /alert/eventhandler/type/{type}/event/{name} [delete]
func DeleteEventHandler(c echo.Context) error {
	eventType := c.Param("type")
	eventHandlerName := c.Param("name")
	err := eventhandler.DeleteEventHandler(eventType, eventHandlerName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, rest.SetMessage(fmt.Sprintf("delete event handler with name %s successfully", eventHandlerName)))
}
