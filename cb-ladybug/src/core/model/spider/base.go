package spider

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cloud-barista/cb-ladybug/src/core/model"
	"github.com/go-resty/resty/v2"

	logger "github.com/sirupsen/logrus"
)

type Model struct {
}

// 결과 처리
func (m *Model) response(resp *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode() > 300 && resp.StatusCode() != http.StatusNotFound {
		logger.Warnf("Spider : statusCode=%d, url=%s, body=%s", resp.StatusCode(), resp.Request.URL, resp)
		status := model.Status{}
		json.Unmarshal(resp.Body(), &status)
		return errors.New(status.Message)
	}
	return nil
}
