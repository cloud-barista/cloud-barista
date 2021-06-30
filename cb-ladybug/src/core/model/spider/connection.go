package spider

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-ladybug/src/utils/app"
	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	logger "github.com/sirupsen/logrus"
)

// connection config
type Connection struct {
	Model
	ConfigName     string `json:"ConfigName"`
	ProviderName   string `json:"ProviderName"`
	DriverName     string `json:"DriverName"`
	CredentialName string `json:"CredentialName"`
	RegionName     string `json:"RegionName"`
}

func NewConnection(name string) *Connection {
	return &Connection{
		Model:      Model{},
		ConfigName: name,
	}
}

// get connection config
func (self *Connection) GET() (bool, error) {

	url := fmt.Sprintf("%s/connectionconfig/%s", *config.Config.SpiderUrl, self.ConfigName)
	resp, err := app.ExecutHTTP(http.MethodGet, url, nil, &self)
	if err != nil {
		return false, err
	}
	if err = self.response(resp, err); err != nil {
		return false, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		logger.Infof("Not found data (status=404, method=%s, url=%s)", http.MethodGet, url)
		return false, nil
	}

	return true, nil

}
