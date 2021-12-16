package template

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	TemplatePath = "/file/templates"
)

type TemplateBuilder struct{}

func RegisterTemplate() {
	rootPath := os.Getenv("CBMON_ROOT")

	// Clean all templates before register
	CleanTemplates()

	// Get all tick template files
	files, err := ioutil.ReadDir(rootPath + TemplatePath)
	if err != nil {
		logrus.Error(fmt.Sprintf("failed to read directory, error=%s", err.Error()))
		return
	}

	// Register tick file from template folder
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		match, _ := regexp.MatchString("^.*\\.(tick$)", f.Name())
		if match {
			// Read tick file
			fileBytes, err := ioutil.ReadFile(rootPath + TemplatePath + "/" + f.Name())
			if err != nil {
				logrus.Errorf("failed to read file, error=%s", err.Error())
				continue
			}

			// Extract file name and content
			tickName := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
			tickScriptContent := string(fileBytes)

			_, err = CreateTemplate(tickName, tickScriptContent)
			if err != nil {
				logrus.Errorf("failed to create template, error=%s", err.Error())
				continue
			}
			logrus.Infof("create tick file with name %s", f.Name())
		}
	}
}

func CleanTemplates() {
	err := DeleteAllTemplates()
	if err != nil {
		logrus.Errorf("failed to clean templates, error=%", err.Error())
	}
}
