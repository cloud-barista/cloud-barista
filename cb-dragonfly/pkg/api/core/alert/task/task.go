package task

import (
	"fmt"
	v1 "github.com/cloud-barista/cb-dragonfly/pkg/storage/metricstore/influxdb/v1"
	"regexp"
	"strings"

	kapacitorclient "github.com/shaodan/kapacitor-client"

	alert "github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/event"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/eventhandler"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/topichandler"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
)

const (
	KapacitorTaskPattern = "dragonfly-*"
	KapacitorTaskFormat  = "dragonfly-%s"
	KapacitorTemplateID  = "default"
	//InfluxDefaultDB      = "cbmon"
	//InfluxDefaultRP      = "autogen"
	AlertMessageFormat = "[{{.Level}}] {{.ID}} {{.TaskName}} Alert \n%s"
)

func ListTasks() ([]types.AlertTask, error) {
	listOpts := kapacitorclient.ListTasksOptions{
		Pattern: KapacitorTaskPattern,
	}
	alertTaskList, err := alert.GetClient().ListTasks(&listOpts)
	if err != nil {
		return nil, err
	}

	alertTaskInfoList := make([]types.AlertTask, len(alertTaskList))
	for idx, alertTask := range alertTaskList {
		alertTaskInfoList[idx] = mappingAlertTaskInfo(alertTask)
	}
	return alertTaskInfoList, nil
}

func GetTask(taskId string) (*types.AlertTask, error) {
	getOpts := kapacitorclient.ListTasksOptions{
		Pattern: fmt.Sprintf(KapacitorTaskFormat, taskId),
	}
	alertTaskList, err := alert.GetClient().ListTasks(&getOpts)
	if err != nil {
		return nil, err
	}
	if len(alertTaskList) == 0 {
		return nil, fmt.Errorf("not found task with ID %s", taskId)
	} else if len(alertTaskList) > 1 {
		return nil, fmt.Errorf("there are multiple tasks with ID %s", taskId)
	}
	alertTask := mappingAlertTaskInfo(alertTaskList[0])
	return &alertTask, nil
}

func CreateTask(alertTaskReq types.AlertTaskReq) (*types.AlertTask, error) {
	createOpts := kapacitorclient.CreateTaskOptions{
		ID:         fmt.Sprintf(KapacitorTaskFormat, alertTaskReq.Name),
		Type:       kapacitorclient.StreamTask,
		TemplateID: fmt.Sprintf(KapacitorTaskFormat, KapacitorTemplateID),
		DBRPs: []kapacitorclient.DBRP{
			{
				Database:        v1.DefaultDatabase,
				RetentionPolicy: v1.CBRetentionPolicyName,
			},
		},
		Status: kapacitorclient.Enabled,
	}
	vars, err := setTemplateVars(alertTaskReq)
	if err != nil {
		return nil, err
	}
	createOpts.Vars = vars

	// Create Alert Task
	alertTask, err := alert.GetClient().CreateTask(createOpts)
	if err != nil {
		return nil, err
	}
	alertTaskInfo := mappingAlertTaskInfo(alertTask)

	// Create Topic Handler
	topicHandlerOpts := map[string]interface{}{}
	if alertTaskReq.AlertEventType == eventhandler.SlackType {
		topicHandlerOpts["workspace"] = alertTaskReq.AlertEventName
	}
	err = topichandler.CreateTopicHandler(fmt.Sprintf(KapacitorTaskFormat, alertTaskInfo.Name), alertTaskInfo.AlertEventType, topicHandlerOpts)
	if err != nil {
		return nil, err
	}
	var dragonflyPort int
	if config.GetInstance().Monitoring.DeployType == "helm" {
		dragonflyPort = config.GetInstance().Dragonfly.HelmPort
	} else {
		dragonflyPort = config.GetInstance().Dragonfly.Port
	}
	// Create Log Topic Handler
	logOpts := map[string]interface{}{
		"url": fmt.Sprintf("http://%s:%d/dragonfly/alert/event", config.GetInstance().Dragonfly.DragonflyIP, dragonflyPort),
	}
	err = topichandler.CreateTopicHandler(fmt.Sprintf(KapacitorTaskFormat, alertTaskInfo.Name), eventhandler.POSTType, logOpts)
	if err != nil {
		return nil, err
	}

	return &alertTaskInfo, nil
}

func UpdateTask(taskId string, alertTaskReq types.AlertTaskReq) (*types.AlertTask, error) {
	taskLink := alert.GetClient().TaskLink(fmt.Sprintf(KapacitorTaskFormat, taskId))
	updateOpts := kapacitorclient.UpdateTaskOptions{}
	vars, err := setTemplateVars(alertTaskReq)
	if err != nil {
		return nil, err
	}
	updateOpts.Vars = vars

	// Update Alert Task
	alertTask, err := alert.GetClient().UpdateTask(taskLink, updateOpts)
	if err != nil {
		return nil, err
	}
	alertTaskInfo := mappingAlertTaskInfo(alertTask)

	// TODO: Update Topic Handler
	topicHandlerOpts := map[string]interface{}{}
	if alertTaskReq.AlertEventType == eventhandler.SlackType {
		topicHandlerOpts["workspace"] = alertTaskReq.AlertEventName
	}
	err = topichandler.UpdateTopicHandler(fmt.Sprintf(KapacitorTaskFormat, alertTaskInfo.Name), alertTaskInfo.AlertEventType, topicHandlerOpts)
	if err != nil {
		return nil, err
	}

	return &alertTaskInfo, nil
}

func DeleteTask(taskId string) error {
	taskLink := alert.GetClient().TaskLink(fmt.Sprintf(KapacitorTaskFormat, taskId))

	alertTask, err := alert.GetClient().Task(taskLink, &kapacitorclient.TaskOptions{})
	if err != nil {
		return err
	}
	alertTaskInfo := mappingAlertTaskInfo(alertTask)

	err = alert.GetClient().DeleteTask(taskLink)
	if err != nil {
		return err
	}

	// Delete Topic
	topicLink := alert.GetClient().TopicLink(fmt.Sprintf(KapacitorTaskFormat, taskId))
	err = alert.GetClient().DeleteTopic(topicLink)
	if err != nil {
		return err
	}

	// Delete Topic Handler
	err = topichandler.DeleteTopicHandler(fmt.Sprintf(KapacitorTaskFormat, taskId), alertTaskInfo.AlertEventType)
	if err != nil {
		return err
	}
	err = topichandler.DeleteTopicHandler(fmt.Sprintf(KapacitorTaskFormat, taskId), eventhandler.POSTType)
	if err != nil {
		return err
	}

	// Delete Event Logs
	err = event.DeleteEventLog(fmt.Sprintf(KapacitorTaskFormat, taskId))
	if err != nil {
		return err
	}
	return nil
}

func setTemplateVars(alertTaskReq types.AlertTaskReq) (map[string]kapacitorclient.Var, error) {
	varMaps := map[string]kapacitorclient.Var{}

	varMaps["measurement"] = newTaskVar(kapacitorclient.VarString, alertTaskReq.Measurement)

	varMaps["target_type"] = newTaskVar(kapacitorclient.VarString, alertTaskReq.TargetType)
	varMaps["target_id"] = newTaskVar(kapacitorclient.VarString, alertTaskReq.TargetId)
	varMaps["where_filter"] = newTaskVar(kapacitorclient.VarLambda, fmt.Sprintf("\"%sId\" == '%s'", strings.ToLower(alertTaskReq.TargetType), alertTaskReq.TargetId))

	varMaps["event_params"] = newTaskVar(kapacitorclient.VarString, alertTaskReq.EventDuration)
	varMaps["event_duration"] = newTaskVar(kapacitorclient.VarDuration, alertTaskReq.EventDuration)
	varMaps["event_interval"] = newTaskVar(kapacitorclient.VarDuration, alertTaskReq.EventDuration)

	varMaps["metric"] = newTaskVar(kapacitorclient.VarString, alertTaskReq.Metric)
	varMaps["alert_math_expression"] = newTaskVar(kapacitorclient.VarString, alertTaskReq.AlertMathExpression)
	varMaps["alert_threshold"] = newTaskVar(kapacitorclient.VarFloat, alertTaskReq.AlertThreshold)
	varMaps["warn_event_cnt"] = newTaskVar(kapacitorclient.VarInt, alertTaskReq.WarnEventCnt)
	varMaps["critic_event_cnt"] = newTaskVar(kapacitorclient.VarInt, alertTaskReq.CriticEventCnt)

	var compareExpression string
	switch alertTaskReq.AlertMathExpression {
	case "equal":
		compareExpression = "=="
	case "greater":
		compareExpression = ">"
	case "equalgreater":
		compareExpression = ">="
	case "less":
		compareExpression = "<"
	case "equalless":
		compareExpression = "<="
	}
	varMaps["state_condition"] = newTaskVar(kapacitorclient.VarLambda, fmt.Sprintf("\"%s\" %s %f", alertTaskReq.Metric, compareExpression, alertTaskReq.AlertThreshold))

	varMaps["warn"] = newTaskVar(kapacitorclient.VarLambda, fmt.Sprintf("\"state_count\" >= %d", alertTaskReq.WarnEventCnt))
	varMaps["crit"] = newTaskVar(kapacitorclient.VarLambda, fmt.Sprintf("\"state_count\" >= %d", alertTaskReq.CriticEventCnt))

	varMaps["alert_event_type"] = newTaskVar(kapacitorclient.VarString, alertTaskReq.AlertEventType)
	varMaps["alert_event_name"] = newTaskVar(kapacitorclient.VarString, alertTaskReq.AlertEventName)
	varMaps["custom_message"] = newTaskVar(kapacitorclient.VarString, alertTaskReq.AlertEventMessage)
	varMaps["alert_message"] = newTaskVar(kapacitorclient.VarString, fmt.Sprintf(AlertMessageFormat, alertTaskReq.AlertEventMessage))
	varMaps["topic_name"] = newTaskVar(kapacitorclient.VarString, fmt.Sprintf(KapacitorTaskFormat, alertTaskReq.Name))

	return varMaps, nil
}

func newTaskVar(varType kapacitorclient.VarType, varVal interface{}) kapacitorclient.Var {
	return kapacitorclient.Var{
		Type:  varType,
		Value: varVal,
	}
}

func mappingAlertTaskInfo(task kapacitorclient.Task) types.AlertTask {
	var taskId string
	uriArr := strings.Split(task.Link.Href, "/")
	if len(uriArr) > 0 {
		taskName := uriArr[len(uriArr)-1]
		reg, _ := regexp.Compile("dragonfly-(.+)")
		if reg.MatchString(taskName) {
			uriParams := reg.FindStringSubmatch(taskName)
			taskId = uriParams[1]
		}
	}

	alertTask := types.AlertTask{
		Name:                taskId,
		Measurement:         getVarByKey(task.Vars, "measurement").(string),
		TargetType:          getVarByKey(task.Vars, "target_type").(string),
		TargetId:            getVarByKey(task.Vars, "target_id").(string),
		EventDuration:       getVarByKey(task.Vars, "event_params").(string),
		Metric:              getVarByKey(task.Vars, "metric").(string),
		AlertMathExpression: getVarByKey(task.Vars, "alert_math_expression").(string),

		AlertThreshold: getVarByKey(task.Vars, "alert_threshold").(float64),
		WarnEventCnt:   getVarByKey(task.Vars, "warn_event_cnt").(int64),
		CriticEventCnt: getVarByKey(task.Vars, "critic_event_cnt").(int64),

		AlertEventType:    getVarByKey(task.Vars, "alert_event_type").(string),
		AlertEventName:    getVarByKey(task.Vars, "alert_event_name").(string),
		AlertEventMessage: getVarByKey(task.Vars, "custom_message").(string),
	}
	return alertTask
}

func getVarByKey(vars kapacitorclient.Vars, key string) interface{} {
	curVar := vars[key]
	return curVar.Value
}
