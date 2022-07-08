package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-mcks/src/core/app"
)

func NewNS(ns string) *NS {
	return &NS{
		Name: ns,
		ID:   ns,
	}
}

func (self *NS) GET() (bool, error) {

	resp, err := executeHTTP(http.MethodGet, fmt.Sprintf("%s/ns/%s", *app.Config.TumblebugUrl, self.Name), nil, &self)
	if err != nil {
		return false, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}
