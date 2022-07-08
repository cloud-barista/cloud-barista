package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-mcks/src/core/app"
)

/* instance of a MCIS */
func NewMCIS(ns string, name string) *MCIS {
	return &MCIS{
		Model: Model{Name: name, Namespace: ns},
		VMs:   []VM{},
	}
}

/* instance of a VM */
func NewVM(namespace string, name string, mcisName string) *VM {
	return &VM{
		Model:       Model{Name: name, Namespace: namespace},
		mcisName:    mcisName,
		UserAccount: VM_USER_ACCOUNT,
	}
}

/* MCIS */
func (self *MCIS) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/mcis/%s", self.Namespace, self.Name), nil, &self)

}

func (self *MCIS) POST() error {

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/mcis", self.Namespace), self, &self)
	if err != nil {
		return err
	}

	return nil
}

func (self *MCIS) DELETE() (bool, error) {

	exist, err := self.GET()
	if err != nil {
		return exist, err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/mcis/%s", self.Namespace, self.Name), nil, app.Status{})
		if err != nil {
			return exist, err
		}
	}

	return exist, nil
}

func (self *MCIS) TERMINATE() error {
	_, err := self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/control/mcis/%s?action=terminate", self.Namespace, self.Name), nil, app.Status{})
	if err != nil {
		return err
	}
	return nil
}

func (self *MCIS) REFINE() error {
	_, err := self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/control/mcis/%s?action=refine", self.Namespace, self.Name), nil, app.Status{})
	if err != nil {
		return err
	}
	return nil
}

/* VM */
func (self *VM) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/mcis/%s/vm/%s", self.Namespace, self.mcisName, self.Name), nil, &self)

}

func (self *VM) POST() error {

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/mcis/%s/vm", self.Namespace, self.mcisName), self, &self)
	if err != nil {
		return err
	}

	return nil

}

func (self *VM) DELETE() (bool, error) {

	exist, err := self.GET()
	if err != nil {
		return exist, err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/mcis/%s/vm/%s", self.Namespace, self.mcisName, self.Name), nil, app.Status{})
		if err != nil {
			return exist, err
		}
	}

	return exist, nil
}
