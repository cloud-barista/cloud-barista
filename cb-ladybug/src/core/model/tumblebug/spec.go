package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	"github.com/go-resty/resty/v2"
)

type Spec struct {
	Model
	Config      string `json:"connectionName"`
	CspSpecName string `json:"cspSpecName"`
	Role        string `json:"role"`
}

func NewSpec(ns string, name string, conf string) *Spec {
	return &Spec{
		Model:  Model{Name: name, namespace: ns},
		Config: conf,
	}
}

func (spec *Spec) GET() (bool, error) {
	// validation
	if err := spec.validate(validation.Validation{}); err != nil {
		return false, err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetBody(fmt.Sprintf(`{"connectionName" : "%s"}`, spec.Config)).
		SetResult(&spec).
		Get(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/spec/%s", spec.namespace, spec.Name))

	if e := spec.response(resp, err); e != nil {
		return false, e
	}
	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}

func (spec *Spec) POST() error {
	// validation
	valid := validation.Validation{}
	valid.Required(spec.CspSpecName, "cspSpecName")
	if err := spec.validate(valid); err != nil {
		return err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetBody(spec).
		SetResult(&spec).
		Post(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/spec", spec.namespace))

	if e := spec.response(resp, err); e != nil {
		return e
	}

	return nil
}

func (spec *Spec) DELETE(ns string) error {
	// validation
	if err := spec.validate(validation.Validation{}); err != nil {
		return err
	}

	exist, err := spec.GET()
	if err != nil {
		return err
	}
	if exist {
		conf := config.Config
		resp, err := resty.New().R().
			SetBasicAuth(conf.Username, conf.Password).
			SetBody(fmt.Sprintf(`{"connectionName" : "%s"}`, spec.Config)).
			SetResult(TumblebugResult{}).
			Delete(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/spec/%s", spec.namespace, spec.Name))

		if e := spec.response(resp, err); e != nil {
			return e
		}
	} else {
		fmt.Println(fmt.Sprintf("delete spec skip.. (name=%s, cause=not found)", spec.Name))
	}

	return nil
}
