package spider

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	"github.com/go-resty/resty/v2"
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
func (region *Region) GET() (bool, error) {

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetResult(&region).
		Get(conf.SpiderUrl + fmt.Sprintf("/region/%s", region.RegionName))

	if err = region.response(resp, err); err != nil {
		return false, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}
