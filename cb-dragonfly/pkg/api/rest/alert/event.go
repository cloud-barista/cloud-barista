package alert

import (
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/event"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/task"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/types"
)

func CreateEventLog(c echo.Context) error {
	var jsonMap map[string]interface{}
	if err := c.Bind(&jsonMap); err != nil {
		return err
	}

	var eventLog types.AlertEventLog
	err := mapstructure.Decode(jsonMap, &eventLog)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = event.CreateEventLog(eventLog)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

// ListEventLog 알람 로그 정보 조회
// @Summary List monitoring alert event
// @Description 알람 로그 정보 목록 조회
// @Tags [Log] Alarm Event Log
// @Accept  json
// @Produce  json
// @Param task_id path string true "태스크 아이디"
// @Success 200 {object} types.AlertEventLog
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /alert/task/{task_id}/events [get]
func ListEventLog(c echo.Context) error {
	taskName := c.Param("task_id")
	logLevel := c.QueryParam("level")
	alertLogList, err := event.ListEventLog(fmt.Sprintf(task.KapacitorTaskFormat, taskName), logLevel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, rest.SetMessage(fmt.Sprintf("failed to get event log list, error=%s", err)))
	}
	return c.JSON(http.StatusOK, alertLogList)
}
