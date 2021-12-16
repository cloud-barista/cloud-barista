package topichandler

import (
	"fmt"

	kapacitorclient "github.com/shaodan/kapacitor-client"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/alert"
)

func CreateTopicHandler(topicName string, eventType string, options map[string]interface{}) error {
	topicLink := alert.GetClient().TopicLink(topicName)

	createOpts := kapacitorclient.TopicHandlerOptions{
		Topic:   topicName,
		ID:      fmt.Sprintf("%s-%s", topicName, eventType),
		Kind:    eventType,
		Options: options,
	}
	_, err := alert.GetClient().CreateTopicHandler(topicLink, createOpts)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTopicHandler(topicName string, eventType string, options map[string]interface{}) error {
	topicHandlerLink := alert.GetClient().TopicHandlerLink(topicName, fmt.Sprintf("%s-%s", topicName, eventType))

	updateOpts := kapacitorclient.TopicHandlerOptions{
		Topic:   topicName,
		ID:      fmt.Sprintf("%s-%s", topicName, eventType),
		Kind:    eventType,
		Options: options,
	}
	_, err := alert.GetClient().ReplaceTopicHandler(topicHandlerLink, updateOpts)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTopicHandler(topicName string, eventType string) error {
	topicHandlerLink := alert.GetClient().TopicHandlerLink(topicName, fmt.Sprintf("%s-%s", topicName, eventType))
	_, err := alert.GetClient().TopicHandler(topicHandlerLink)
	if err != nil {
		return err
	}
	err = alert.GetClient().DeleteTopicHandler(topicHandlerLink)
	if err != nil {
		return err
	}
	return nil
}
