package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	"github.com/go-resty/resty/v2"
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

func (ssh *SSHKey) GET() (bool, error) {
	// validation
	if err := ssh.validate(validation.Validation{}); err != nil {
		return false, err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetBody(fmt.Sprintf(`{"connectionName" : "%s"}`, ssh.Config)).
		SetResult(&ssh).
		Get(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/sshKey/%s", ssh.namespace, ssh.Name))

	if e := ssh.response(resp, err); e != nil {
		return false, e
	}
	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}

func (ssh *SSHKey) POST() error {
	// validation
	valid := validation.Validation{}
	valid.Required(ssh.Username, "username")
	if err := ssh.validate(valid); err != nil {
		return err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetBody(ssh).
		SetResult(&ssh).
		Post(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/sshKey", ssh.namespace))

	if err = ssh.response(resp, err); err != nil {
		return err
	}

	return nil
}

func (ssh *SSHKey) DELETE(ns string) error {
	// validation
	if err := ssh.validate(validation.Validation{}); err != nil {
		return err
	}

	exist, err := ssh.GET()
	if err != nil {
		return err
	}
	if exist {
		conf := config.Config

		resp, err := resty.New().R().
			SetBasicAuth(conf.Username, conf.Password).
			SetBody(fmt.Sprintf(`{"connectionName" : "%s"}`, ssh.Config)).
			SetResult(TumblebugResult{}).
			Delete(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/sshKey/%s", ssh.namespace, ssh.Name))

		if err = ssh.response(resp, err); err != nil {
			return err
		}
	} else {
		fmt.Println(fmt.Sprintf("delete sshkey skip (name=%s, cause=not found)", ssh.Name))
	}

	return nil
}
