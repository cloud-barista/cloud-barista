package tumblebug

import (
	"fmt"
	"net/http"
)

type Region struct {
	Model
	RegionName       string     `json:"RegionName"`
	ProviderName     string     `json:"ProviderName"`
	KeyValueInfoList []KeyValue `json:"KeyValueInfoList"`
}

func NewRegion(name string) *Region {
	return &Region{
		Model:      Model{},
		RegionName: name,
	}
}

// get region
func (self *Region) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/region/%s", self.RegionName), nil, &self)

}
