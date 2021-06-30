package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-ladybug/src/utils/app"
	"github.com/cloud-barista/cb-ladybug/src/utils/config"
)

type NS struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewNS(ns string) *NS {
	return &NS{
		Name: ns,
		ID:   ns,
	}
}

func (self *NS) GET() (bool, error) {

	resp, err := app.ExecutHTTP(http.MethodGet, fmt.Sprintf("%s/ns/%s", *config.Config.TumblebugUrl, self.Name), nil, &self)
	if err != nil {
		return false, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}
