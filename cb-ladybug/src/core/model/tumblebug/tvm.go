package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/cloud-barista/cb-ladybug/src/core/model"
	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	"github.com/go-resty/resty/v2"
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

func (tvm *TVM) GET() (bool, error) {
	// validation
	if err := tvm.validate(validation.Validation{}); err != nil {
		return false, err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetResult(&tvm.VM).
		Get(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/mcis/%s/vm/%s", tvm.namespace, tvm.Name, tvm.VM.Name))

	if e := tvm.response(resp, err); e != nil {
		return false, e
	}
	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}

func (tvm *TVM) POST() error {
	// validation
	if err := tvm.validate(validation.Validation{}); err != nil {
		return err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetBody(tvm.VM).
		SetResult(&tvm.VM).
		Post(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/mcis/%s/vm", tvm.namespace, tvm.Name))

	if e := tvm.response(resp, err); e != nil {
		return e
	}

	return nil
}

func (tvm *TVM) DELETE() error {
	// validation
	if err := tvm.validate(validation.Validation{}); err != nil {
		return err
	}

	exist, err := tvm.GET()
	if err != nil {
		return err
	}
	if exist {
		conf := config.Config
		resp, err := resty.New().R().
			SetBasicAuth(conf.Username, conf.Password).
			SetResult(TumblebugResult{}).
			Delete(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/mcis/%s/vm/%s", tvm.namespace, tvm.Name, tvm.VM.Name))

		if err = tvm.response(resp, err); err != nil {
			return err
		}
	} else {
		fmt.Println(fmt.Sprintf("delete VM skip (name=%s, cause=not found)", tvm.VM.Name))
	}

	return nil
}
