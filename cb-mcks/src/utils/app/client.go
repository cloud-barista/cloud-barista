package app

import (
	"github.com/cloud-barista/cb-mcks/src/utils/config"
	"github.com/go-resty/resty/v2"
)

func ExecuteHTTP(method string, url string, body interface{}, result interface{}) (*resty.Response, error) {

	conf := config.Config

	req := resty.New().SetDisableWarn(true).R().SetBasicAuth(*conf.Username, *conf.Password)

	if body != nil {
		req.SetBody(body)
	}
	if result != nil {
		req.SetResult(result)
	}

	// execute
	return req.Execute(method, url)

}
