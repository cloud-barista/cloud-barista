package spider

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-ladybug/src/utils/app"
	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	logger "github.com/sirupsen/logrus"
)

type Region struct {
	Model
	RegionName       string     `json:"RegionName"`
	ProviderName     string     `json:"ProviderName"`
	KeyValueInfoList []KeyValue `json:"KeyValueInfoList"`
}

type KeyValue struct {
	Key   string
	Value string
}

func NewRegion(name string) *Region {
	return &Region{
		Model:      Model{},
		RegionName: name,
	}
}

// get region
func (self *Region) GET() (bool, error) {

	url := fmt.Sprintf("%s/region/%s", *config.Config.SpiderUrl, self.RegionName)
	resp, err := app.ExecutHTTP(http.MethodGet, url, nil, &self)
	if err != nil {
		return false, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		logger.Warnf("Not found data (status=404, method=%s, url=%s)", http.MethodGet, url)
		return false, nil
	}

	return true, nil

}
