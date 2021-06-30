package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-ladybug/src/core/model"

	logger "github.com/sirupsen/logrus"
)

type MCIS struct {
	Model
	Description string     `json:"description"`
	VMs         []model.VM `json:"vm"` // output
}

func NewMCIS(ns string, name string) *MCIS {
	return &MCIS{
		Model: Model{Name: name, namespace: ns},
		VMs:   []model.VM{},
	}
}

func (self *MCIS) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/mcis/%s", self.namespace, self.Name), nil, &self)

	// if err = self.hasResponseMessage(resp); err != nil {
	// 	// 이부분은 향후 tumblebug 에서 수정해야 할듯
	// 	if strings.HasPrefix(err.Error(), "Cannot find") {
	// 		return false, nil
	// 	} else {
	// 		return false, err
	// 	}
	// }
}

func (self *MCIS) POST() error {

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/mcis", self.namespace), self, &self)
	if err != nil {
		return err
	}

	return nil
}

func (self *MCIS) DELETE() error {

	exist, err := self.GET()
	if err != nil {
		return err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/mcis/%s", self.namespace, self.Name), nil, model.Status{})
		if err != nil {
			return err
		}
	} else {
		logger.Infof("delete MCIS skip (name=%s, cause=not found)", self.Name)
	}

	return nil
}

func (self *MCIS) TERMINATE() error {
	_, err := self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/mcis/%s?action=terminate", self.namespace, self.Name), nil, model.Status{})
	if err != nil {
		return err
	}
	return nil
}

func (self *MCIS) REFINE() error {
	_, err := self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/mcis/%s?action=refine", self.namespace, self.Name), nil, model.Status{})
	if err != nil {
		return err
	}
	return nil
}
