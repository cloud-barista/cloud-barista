package tumblebug

import (
	"fmt"
	"net/http"
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

	return self.execute(http.MethodGet, fmt.Sprintf("/connConfig/%s", self.ConfigName), nil, &self)

}
