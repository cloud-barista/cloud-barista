package template

import (
	"fmt"

	kapacitorclient "github.com/shaodan/kapacitor-client"
	"github.com/sirupsen/logrus"

	"github.com/cloud-barista/cb-dragonfly/pkg/core/alert"
)

const (
	NamePattern   = "dragonfly-*"
	FormatPattern = "dragonfly-%s"
)

func ListTemplates() ([]kapacitorclient.Template, error) {
	listOpts := kapacitorclient.ListTemplatesOptions{
		Pattern: NamePattern,
	}
	return alert.GetClient().ListTemplates(&listOpts)
}

func GetTemplate(templateName string) (*kapacitorclient.Template, error) {
	templateLink, err := getTemplateLinkByName(templateName)
	if err != nil {
		return nil, fmt.Errorf("not found template with Name %s", templateName)
	}
	template, err := alert.GetClient().Template(*templateLink, &kapacitorclient.TemplateOptions{})
	if err != nil {
		return nil, fmt.Errorf("not found template with Name %s", templateName)
	}
	return &template, nil
}

func CreateTemplate(templateName string, tickScript string) (*kapacitorclient.Template, error) {
	createOpts := kapacitorclient.CreateTemplateOptions{
		ID:         fmt.Sprintf(FormatPattern, templateName),
		Type:       kapacitorclient.StreamTask,
		TICKscript: tickScript,
	}
	alertTemplate, err := alert.GetClient().CreateTemplate(createOpts)
	if err != nil {
		return nil, err
	}
	return &alertTemplate, nil
}

func UpdateTemplate(templateId string) (*kapacitorclient.Template, error) {
	// TODO: implement method
	return nil, nil
}

func DeleteTemplate(templateName string) error {
	templateLink, err := getTemplateLinkByName(templateName)
	if err != nil {
		return fmt.Errorf("not found template with Name %s", templateName)
	}
	err = alert.GetClient().DeleteTask(*templateLink)
	if err != nil {
		return fmt.Errorf("failed to delete template, err=%s", err.Error())
	}
	return nil
}

func DeleteAllTemplates() error {
	templateList, err := ListTemplates()
	if err != nil {
		return fmt.Errorf("failed to get list of templates, err=%s", err.Error())
	}
	for _, tmpl := range templateList {
		err = alert.GetClient().DeleteTask(tmpl.Link)
		if err != nil {
			logrus.Errorf("failed to delete template with name %s, error=%s", tmpl.ID, err.Error())
			continue
		}
	}
	return nil
}

func getTemplateLinkByName(templateName string) (*kapacitorclient.Link, error) {
	getOpts := kapacitorclient.ListTemplatesOptions{
		Pattern: fmt.Sprintf(FormatPattern, templateName),
	}
	alertTemplateList, err := alert.GetClient().ListTemplates(&getOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list template, err=%s", err.Error())
	}
	if len(alertTemplateList) == 0 {
		return nil, fmt.Errorf("not found template with Name %s", templateName)
	} else if len(alertTemplateList) > 1 {
		return nil, fmt.Errorf("there are multiple templates with Name %s", templateName)
	}
	return &alertTemplateList[0].Link, nil
}
