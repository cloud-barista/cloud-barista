package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/beego/beego/v2/core/validation"
	"github.com/cloud-barista/cb-mcks/src/core/model"

	logger "github.com/sirupsen/logrus"
)

type Spec struct {
	Model
	Config      string `json:"connectionName"`
	CspSpecName string `json:"cspSpecName"`
}

func NewSpec(ns string, name string, conf string) *Spec {
	return &Spec{
		Model:  Model{Name: name, namespace: ns},
		Config: conf,
	}
}

func (self *Spec) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/resources/spec/%s", self.namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), &self)

}

func (self *Spec) POST() error {

	// validation
	valid := validation.Validation{}
	valid.Required(self.CspSpecName, "cspSpecName")
	if err := self.validate(valid); err != nil {
		return err
	}

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/resources/spec", self.namespace), self, &self)
	if err != nil {
		return err
	}

	return nil

}

func (self *Spec) DELETE(ns string) error {

	exist, err := self.GET()
	if err != nil {
		return err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/resources/spec/%s", self.namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), model.Status{})
		if err != nil {
			return err
		}

	} else {
		logger.Infof("delete spec skip.. (name=%s, cause=not found)", self.Name)
	}

	return nil
}
