package spider

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	"github.com/go-resty/resty/v2"
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
func (conn *Connection) GET() (bool, error) {

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetResult(&conn).
		Get(conf.SpiderUrl + fmt.Sprintf("/connectionconfig/%s", conn.ConfigName))

	if err = conn.response(resp, err); err != nil {
		return false, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}
