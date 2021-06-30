package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-ladybug/src/core/model"

	logger "github.com/sirupsen/logrus"
)

type TVM struct {
	Model
	VM model.VM
}

func NewTVm(ns string, mcisName string) *TVM {
	return &TVM{
		Model: Model{namespace: ns, Name: mcisName},
	}
}

func (self *TVM) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/mcis/%s/vm/%s", self.namespace, self.Name, self.VM.Name), nil, &self.VM)

}

func (self *TVM) POST() error {

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/mcis/%s/vm", self.namespace, self.Name), self.VM, &self.VM)
	if err != nil {
		return err
	}

	return nil

}

func (self *TVM) DELETE() error {

	exist, err := self.GET()
	if err != nil {
		return err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/mcis/%s/vm/%s", self.namespace, self.Name, self.VM.Name), nil, model.Status{})
		if err != nil {
			return err
		}
	} else {
		logger.Infof("delete VM skip (name=%s, cause=not found)", self.VM.Name)
	}

	return nil
}
