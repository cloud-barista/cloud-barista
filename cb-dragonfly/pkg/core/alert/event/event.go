package event

import (
	"encoding/json"

	"github.com/cloud-barista/cb-dragonfly/pkg/core/alert/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/localstore"
)

func CreateEventLog(eventLog types.AlertEventLog) error {
	var eventLogArr []types.AlertEventLog

	eventLogStr := localstore.GetInstance().StoreGet(eventLog.Id)

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
	err = localstore.GetInstance().StorePut(eventLog.Id, string(newEventLogBytes))
	if err != nil {
		return err
	}
	return nil
}

func ListEventLog(taskId string, logLevel string) ([]types.AlertEventLog, error) {
	eventLogArr := []types.AlertEventLog{}
	eventLogStr := localstore.GetInstance().StoreGet(taskId)
	//if err != nil {
	//	return nil, err
	//}
	//if eventLogStr == nil {
	//	return eventLogArr, nil
	//}
	err := json.Unmarshal([]byte(eventLogStr), &eventLogArr)
	if err != nil {
		return nil, err
	}
	return eventLogArr, nil
}

func DeleteEventLog(taskId string) error {
	return localstore.GetInstance().StoreDelete(taskId)
}
