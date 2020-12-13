package tumblebug

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/cloud-barista/cb-ladybug/src/core/model"
	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	"github.com/go-resty/resty/v2"
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

func (mcis *MCIS) GET() (bool, error) {
	// validation
	if err := mcis.validate(validation.Validation{}); err != nil {
		return false, err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetResult(&mcis).
		Get(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/mcis/%s", mcis.namespace, mcis.Name))

	if err = mcis.response(resp, err); err != nil {
		return false, err
	}
	// 존재여부 확인
	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}
	if err = mcis.hasResponseMessage(resp); err != nil {
		// 이부분은 향후 tumblebug 에서 수정해야 할듯
		if strings.HasPrefix(err.Error(), "Cannot find") {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func (mcis *MCIS) POST() error {
	// validation
	if err := mcis.validate(validation.Validation{}); err != nil {
		return err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetBody(mcis).
		SetResult(&mcis).
		Post(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/mcis", mcis.namespace))

	if err = mcis.response(resp, err); err != nil {
		return err
	}

	return nil
}

func (mcis *MCIS) DELETE() error {
	// validation
	if err := mcis.validate(validation.Validation{}); err != nil {
		return err
	}

	exist, err := mcis.GET()
	if err != nil {
		return err
	}
	if exist {
		conf := config.Config
		resp, err := resty.New().R().
			SetBasicAuth(conf.Username, conf.Password).
			SetResult(TumblebugResult{}).
			Delete(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/mcis/%s", mcis.namespace, mcis.Name))

		if err = mcis.response(resp, err); err != nil {
			return err
		}
	} else {
		fmt.Println(fmt.Sprintf("delete MCIS skip (name=%s, cause=not found)", mcis.Name))
	}

	return nil
}
