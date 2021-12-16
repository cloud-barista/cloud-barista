package slack

import (
	"fmt"
	"regexp"

	kapacitorclient "github.com/shaodan/kapacitor-client"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert/types"
)

const (
	EventType = "slack"
)

type SlackHandler struct{}

func (s SlackHandler) ListEventHandlers() ([]types.AlertEventHandler, error) {
	configSections, err := alert.GetClient().ConfigSections()
	if err != nil {
		return nil, err
	}

	slackConfigSection, ok := configSections.Sections[EventType]
	if !ok {
		return nil, fmt.Errorf("not found event type with Name %s", EventType)
	}

	var eventHandlerList []types.AlertEventHandler
	for _, configElement := range slackConfigSection.Elements {
		eventHandlerInfo := mappingAlertEventHandlerInfo(configElement)
		if eventHandlerInfo == nil {
			continue
		}
		eventHandlerList = append(eventHandlerList, *eventHandlerInfo)
	}
	return eventHandlerList, nil
}

func (s SlackHandler) GetEventHandler(name string) (types.AlertEventHandler, error) {
	slackLink := alert.GetClient().ConfigElementLink(EventType, name)
	configElement, err := alert.GetClient().ConfigElement(slackLink)
	if err != nil {
		return types.AlertEventHandler{}, fmt.Errorf("not found event handler with Name %s", name)
	}
	eventHandlerInfo := mappingAlertEventHandlerInfo(configElement)
	return *eventHandlerInfo, nil
}

func (s SlackHandler) CreateEventHandler(createOpts types.AlertEventHandlerReq) (types.AlertEventHandler, error) {
	sectionLink := alert.GetClient().ConfigSectionLink(EventType)

	// Set slack create options
	options := map[string]interface{}{}
	options["enabled"] = true
	options["workspace"] = createOpts.Name
	options["url"] = createOpts.Url
	options["channel"] = createOpts.Channel

	// Create slack event handler
	err := alert.GetClient().ConfigUpdate(sectionLink, kapacitorclient.ConfigUpdateAction{
		Add: options,
	})
	if err != nil {
		return types.AlertEventHandler{}, err
	}
	return s.GetEventHandler(createOpts.Name)
}

func (s SlackHandler) UpdateEventHandler(name string, updateOpts types.AlertEventHandlerReq) (types.AlertEventHandler, error) {
	slackLink := alert.GetClient().ConfigElementLink(EventType, name)

	// Set slack update options
	options := map[string]interface{}{}
	options["enabled"] = true
	options["workspace"] = updateOpts.Name
	options["url"] = updateOpts.Url
	options["channel"] = updateOpts.Channel

	// Update slack event handler
	err := alert.GetClient().ConfigUpdate(slackLink, kapacitorclient.ConfigUpdateAction{
		Set: options,
	})
	if err != nil {
		return types.AlertEventHandler{}, err
	}
	return s.GetEventHandler(name)
}

func (s SlackHandler) DeleteEventHandler(name string) error {
	sectionLink := alert.GetClient().ConfigSectionLink(EventType)

	// Delete slack event handler
	err := alert.GetClient().ConfigUpdate(sectionLink, kapacitorclient.ConfigUpdateAction{
		Remove: []string{name},
	})
	if err != nil {
		return err
	}
	return nil
}

func mappingAlertEventHandlerInfo(configElement kapacitorclient.ConfigElement) *types.AlertEventHandler {
	reg, _ := regexp.Compile(fmt.Sprintf("/kapacitor/v1/config/%s/(.+)", EventType))
	if !reg.MatchString(configElement.Link.Href) {
		return nil
	}
	uriParams := reg.FindStringSubmatch(configElement.Link.Href)
	alertEventHandler := types.AlertEventHandler{
		ID:   configElement.Link.Href,
		Type: EventType,
		Name: uriParams[1],
		Options: map[string]interface{}{
			"url":     configElement.Options["url"],
			"channel": configElement.Options["channel"],
		},
	}
	return &alertEventHandler
}
