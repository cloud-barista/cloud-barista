package spider

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-ladybug/src/core/model"
	"github.com/go-resty/resty/v2"
)

type Model struct {
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
