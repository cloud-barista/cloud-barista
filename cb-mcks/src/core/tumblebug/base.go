package tumblebug

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/core/validation"
	"github.com/cloud-barista/cb-mcks/src/core/app"
	"github.com/go-resty/resty/v2"

	logger "github.com/sirupsen/logrus"
)

/* execute  */
func (self *Model) execute(method string, url string, body interface{}, result interface{}) (bool, error) {

	if strings.Contains(url, "/ns") {
		// validation
		if err := self.validate(validation.Validation{}); err != nil {
			return false, err
		}
	}

	logger.Debugf("[%s] Start to execute a HTTP (url='%s')", method, url)
	resp, err := executeHTTP(method, *app.Config.TumblebugUrl+url, body, result)
	if err != nil {
		return false, err
	}

	// response check
	if resp.StatusCode() == http.StatusNotFound {
		logger.Infof("[%s] Could not be found data. (method=%s, url='%s')", self.Name, method, url)
		return false, nil
	} else if resp.StatusCode() > 300 {
		logger.Warnf("[%s] Received error data from the Tumblebug (statusCode=%d, url='%s', body='%v')", self.Name, resp.StatusCode(), resp.Request.URL, resp)
		status := app.Status{}
		json.Unmarshal(resp.Body(), &status)
		// message > message 로 리턴되는 경우가 있어서 한번더 unmarshal 작업
		if json.Valid([]byte(status.Message)) {
			json.Unmarshal([]byte(status.Message), &status)
		}
		return false, errors.New(status.Message)
	}

	return true, nil

}

/* execute HTTP */
func executeHTTP(method string, url string, body interface{}, result interface{}) (*resty.Response, error) {

	req := resty.New().SetDisableWarn(true).R().SetBasicAuth(*app.Config.Username, *app.Config.Password)

	if body != nil {
		req.SetBody(body)
	}
	if result != nil {
		req.SetResult(result)
	}

	// execute
	return req.Execute(method, url)

}

/* validate a namespace & a name */
func (self *Model) validate(valid validation.Validation) error {
	valid.Required(self.Namespace, "namespace")
	valid.Required(self.Name, "name")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return errors.New(fmt.Sprintf("Invalid field '%s'%s.", err.Key, err.Error()))
		}
	}
	return nil
}
