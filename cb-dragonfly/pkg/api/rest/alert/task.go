package alert

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	"github.com/cloud-barista/cb-dragonfly/pkg/core/alert/task"
	"github.com/cloud-barista/cb-dragonfly/pkg/core/alert/types"
)

// 모니터링 알람 목록 조회
func ListAlertTask(c echo.Context) error {
	alertTaskList, err := task.ListTasks()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, alertTaskList)
}

// 모니터링 알람 조회
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

// 모니터링 알람 생성
func CreateAlertTask(c echo.Context) error {
	createTaskReq, err := setAlertTaskReq(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	alertTask, err := task.CreateTask(*createTaskReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, *alertTask)
}

// 모니터링 수정 생성
func UpdateAlertTask(c echo.Context) error {
	updateTaskReq, err := setAlertTaskReq(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	alertTask, err := task.UpdateTask(updateTaskReq.Name, *updateTaskReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, *alertTask)
}

// 모니터링 알람 삭제
func DeleteAlertTask(c echo.Context) error {
	taskId := c.Param("task_id")
	err := task.DeleteTask(taskId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, rest.SetMessage(fmt.Sprintf("delete alert task with name %s successfully", taskId)))
}

func setAlertTaskReq(c echo.Context) (*types.AlertTaskReq, error) {
	alertTaskReq := &types.AlertTaskReq{
		Name:                c.FormValue("name"),
		Measurement:         c.FormValue("measurement"),
		TargetType:          c.FormValue("target_type"),
		TargetId:            c.FormValue("target_id"),
		EventDuration:       c.FormValue("event_duration"),
		Metric:              c.FormValue("metric"),
		AlertMathExpression: c.FormValue("alert_math_expression"),
		AlertEventType:      c.FormValue("alert_event_type"),
		AlertEventName:      c.FormValue("alert_event_name"),
		AlertEventMessage:   c.FormValue("alert_event_message"),
	}

	if alertThreshold, err := strconv.ParseFloat(c.FormValue("alert_threshold"), 64); err != nil {
		return nil, err
	} else {
		alertTaskReq.AlertThreshold = alertThreshold
	}

	if warnEventCnt, err := strconv.ParseInt(c.FormValue("warn_event_cnt"), 10, 64); err != nil {
		return nil, err
	} else {
		alertTaskReq.WarnEventCnt = warnEventCnt
	}
	if criticEventCnt, err := strconv.ParseInt(c.FormValue("critic_event_cnt"), 10, 64); err != nil {
		return nil, err
	} else {
		alertTaskReq.CriticEventCnt = criticEventCnt
	}
	return alertTaskReq, nil
}
