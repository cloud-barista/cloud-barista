package smtp

import (
	"errors"
	"fmt"

	kapacitorclient "github.com/shaodan/kapacitor-client"

	"github.com/cloud-barista/cb-dragonfly/pkg/core/alert"
	"github.com/cloud-barista/cb-dragonfly/pkg/core/alert/types"
)

const (
	EventType = "smtp"
)

type SmtpHandler struct{}

func (s SmtpHandler) ListEventHandlers() ([]types.AlertEventHandler, error) {
	eventHandlerInfo, err := s.GetEventHandler("")
	if err != nil {
		return nil, err
	}
	return []types.AlertEventHandler{eventHandlerInfo}, nil
}

func (s SmtpHandler) GetEventHandler(name string) (types.AlertEventHandler, error) {
	smtpLink := alert.GetClient().ConfigSectionLink(EventType)
	smtpConfigSection, err := alert.GetClient().ConfigSection(smtpLink)
	if err != nil {
		return types.AlertEventHandler{}, err
	}
	if len(smtpConfigSection.Elements) == 1 {
		eventHandlerInfo := mappingAlertEventHandlerInfo(smtpConfigSection.Elements[0])
		return *eventHandlerInfo, nil
	}
	return types.AlertEventHandler{}, errors.New("failed to get smtp event handler")
}

func (s SmtpHandler) CreateEventHandler(createOpts types.AlertEventHandlerReq) (types.AlertEventHandler, error) {
	return types.AlertEventHandler{}, errors.New("SMTP event handler can not create new event handler")
}

func (s SmtpHandler) UpdateEventHandler(name string, updateOpts types.AlertEventHandlerReq) (types.AlertEventHandler, error) {
	defaultSmtpLink := kapacitorclient.Link{
		Relation: kapacitorclient.Self,
		Href:     fmt.Sprintf("/kapacitor/v1/config/%s/", EventType),
	}

	// Set smtp create options
	options := map[string]interface{}{}
	options["enabled"] = true
	options["host"] = updateOpts.Host
	options["port"] = updateOpts.Port
	options["from"] = updateOpts.From
	options["to"] = updateOpts.To
	options["username"] = updateOpts.Username
	options["password"] = updateOpts.Password

	// Create smtp event handler
	err := alert.GetClient().ConfigUpdate(defaultSmtpLink, kapacitorclient.ConfigUpdateAction{
		Set: options,
	})
	if err != nil {
		return types.AlertEventHandler{}, err
	}
	return s.GetEventHandler(updateOpts.Name)
}

func (s SmtpHandler) DeleteEventHandler(name string) error {
	return errors.New("SMTP event handler can not delete default event handler")
}

func mappingAlertEventHandlerInfo(configElement kapacitorclient.ConfigElement) *types.AlertEventHandler {
	alertEventHandler := types.AlertEventHandler{
		ID:   configElement.Link.Href,
		Type: EventType,
		Name: EventType,
		Options: map[string]interface{}{
			"host":     configElement.Options["host"],
			"port":     configElement.Options["port"],
			"from":     configElement.Options["from"],
			"to":       configElement.Options["to"],
			"username": configElement.Options["username"],
			"password": configElement.Options["password"],
		},
	}
	return &alertEventHandler
}
