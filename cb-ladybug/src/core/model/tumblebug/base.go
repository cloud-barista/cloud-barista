package tumblebug

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/cloud-barista/cb-ladybug/src/core/model"
	"github.com/go-resty/resty/v2"
)

type KeyValue struct {
	Key   string
	Value string
}

type Model struct {
	Name      string `json:"name"`
	namespace string
}

// 결과 처리
func (m *Model) response(resp *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode() > 300 && resp.StatusCode() != http.StatusNotFound {
		fmt.Println(fmt.Sprintf("[Error Response] %d, url=%s, resp=%s", resp.StatusCode(), resp.Request.URL, resp))
		status := model.Status{}
		json.Unmarshal(resp.Body(), &status)
		return errors.New(status.Message)
	}
	return nil
}

func (m *Model) validate(valid validation.Validation) error {
	valid.Required(m.namespace, "namespace")
	valid.Required(m.Name, "name")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(fmt.Sprintf("[%s]%s", err.Key, err.Error()))
		}
	}
	return nil
}

// status :200 , body = {message: "Cannot find ..." }  형태의 response 에러처리
func (m *Model) hasResponseMessage(resp *resty.Response) error {
	var d map[string]interface{}
	json.Unmarshal(resp.Body(), &d)
	if d["message"] != nil {
		return errors.New(fmt.Sprintf("%s", d["message"]))
	}
	return nil
}
