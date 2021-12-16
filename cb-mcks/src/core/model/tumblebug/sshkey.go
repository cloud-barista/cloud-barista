package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/beego/beego/v2/core/validation"
	"github.com/cloud-barista/cb-mcks/src/core/model"

	logger "github.com/sirupsen/logrus"
)

type SSHKey struct {
	Model
	Config     string `json:"connectionName"`
	Username   string `json:"username"`
	PrivateKey string `json:"privateKey"` // output
}

func NewSSHKey(ns string, name string, conf string) *SSHKey {
	return &SSHKey{
		Model:  Model{Name: name, namespace: ns},
		Config: conf,
	}
}

func (self *SSHKey) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/resources/sshKey/%s", self.namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), &self)

}

func (self *SSHKey) POST() error {
	// validation
	valid := validation.Validation{}
	valid.Required(self.Username, "username")
	if err := self.validate(valid); err != nil {
		return err
	}

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/resources/sshKey", self.namespace), self, &self)
	if err != nil {
		return err
	}

	return nil

}

func (self *SSHKey) DELETE(ns string) error {

	exist, err := self.GET()
	if err != nil {
		return err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/resources/sshKey/%s", self.namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), model.Status{})
		if err != nil {
			return err
		}
	} else {
		logger.Infof("delete sshkey skip (name=%s, cause=not found)", self.Name)
	}

	return nil
}
