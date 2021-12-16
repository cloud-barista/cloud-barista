package event

import (
	"encoding/json"
	"fmt"
	"strings"

	alerttypes "github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/storage/cbstore"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
)

func CreateEventLog(eventLog alerttypes.AlertEventLog) error {
	var eventLogArr []alerttypes.AlertEventLog

	eventLogStr := cbstore.GetInstance().StoreGet(eventLog.Id)

	if eventLogStr != "" {
		// Get event log array
		err := json.Unmarshal([]byte(eventLogStr), &eventLogArr)
		if err != nil {
			return err
		}
	}

	// Add new event log
	eventLogArr = append(eventLogArr, eventLog)

	// Save event log
	newEventLogBytes, err := json.Marshal(eventLogArr)
	if err != nil {
		return err
	}
	err = cbstore.GetInstance().StorePut(fmt.Sprintf("%s/%s", types.EventLog, eventLog.Id), string(newEventLogBytes))
	if err != nil {
		return err
	}
	return nil
}

func ListEventLog(taskId string, logLevel string) ([]alerttypes.AlertEventLog, error) {
	var eventLogArr []alerttypes.AlertEventLog
	eventLogStr := cbstore.GetInstance().StoreGet(taskId)
	if eventLogStr == "" {
		return []alerttypes.AlertEventLog{}, nil
	}
	err := json.Unmarshal([]byte(eventLogStr), &eventLogArr)
	if err != nil {
		return nil, err
	}
	if logLevel == "" {
		return eventLogArr, nil
	}

	filterdEventLogArr := []alerttypes.AlertEventLog{}
	for _, log := range eventLogArr {
		if strings.EqualFold(log.Level, logLevel) {
			filterdEventLogArr = append(filterdEventLogArr, log)
		}
	}
	return filterdEventLogArr, nil
}

func DeleteEventLog(taskId string) error {
	return cbstore.GetInstance().StoreDelete(taskId)
}
