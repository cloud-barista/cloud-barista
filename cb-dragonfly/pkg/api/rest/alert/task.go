package alert

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/task"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
)

// ListAlertTask 알람 목록 조회
// @Summary List monitoring alert
// @Description 알람 목록 조회
// @Tags [Alarm] Alarm management
// @Accept  json
// @Produce  json
// @Success 200 {object} []types.AlertTask
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /alert/tasks [get]
func ListAlertTask(c echo.Context) error {
	alertTaskList, err := task.ListTasks()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, alertTaskList)
}

// GetAlertTask 알람 조회
// @Summary Get monitoring alert
// @Description 알람 조회
// @Tags [Alarm] Alarm management
// @Accept  json
// @Produce  json
// @Param task_id path string true "태스크 아이디"
// @Success 200 {object} types.AlertTask
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /alert/task/{task_id} [get]
func GetAlertTask(c echo.Context) error {
	taskId := c.Param("task_id")
	alertTask, err := task.GetTask(taskId)
	if err != nil {
		if strings.Contains(strings.ToUpper(err.Error()), "NOT FOUND") {
			return echo.NewHTTPError(http.StatusNotFound, rest.SetMessage(err.Error()))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, *alertTask)
}

// CreateAlertTask 알람 생성
// @Summary Create monitoring alert
// @Description 알람 생성
// @Tags [Alarm] Alarm management
// @Accept  json
// @Produce  json
// @Param eventHandlerInfo body types.AlertTask true "Details for an Event object"
// @Success 200 {object} types.AlertTask
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /alert/task [post]
func CreateAlertTask(c echo.Context) error {
	params := &types.AlertTaskReq{}
	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	alertTask, err := task.CreateTask(*params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, *alertTask)
}

// UpdateAlertTask 알람 수정
// @Summary Update monitoring alert
// @Description 알람 수정
// @Tags [Alarm] Alarm management
// @Accept  json
// @Produce  json
// @Param task_id path string true "태스크 아이디"
// @Param eventHandlerInfo body types.AlertTask true "Details for an Event object"
// @Success 200 {object} types.AlertTask
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /alert/task/{task_id} [put]
func UpdateAlertTask(c echo.Context) error {
	params := &types.AlertTaskReq{}
	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	alertTask, err := task.UpdateTask(params.Name, *params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, *alertTask)
}

// DeleteAlertTask 알람 삭제
// @Summary Delete monitoring alert
// @Description 알람 삭제
// @Tags [Alarm] Alarm management
// @Accept  json
// @Produce  json
// @Param task_id path string true "태스크 아이디"
// @Success 200 {object} rest.SimpleMsg
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /alert/task/{task_id} [delete]
func DeleteAlertTask(c echo.Context) error {
	taskId := c.Param("task_id")
	err := task.DeleteTask(taskId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, rest.SetMessage(fmt.Sprintf("delete alert task with name %s successfully", taskId)))
}
