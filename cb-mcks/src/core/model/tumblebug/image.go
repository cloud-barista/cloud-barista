package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-mcks/src/core/model"
	logger "github.com/sirupsen/logrus"
)

type Image struct {
	Model
	Config       string     `json:"connectionName"`
	CspImageId   string     `json:"cspImageId"`
	CspImageName string     `json:"cspImageName"` // output
	CreationDate string     `json:"creationDate"` // output
	Description  string     `json:"description"`  //
	GuestOS      string     `json:"guestOS"`      //
	Status       string     `json:"status"`       // output
	KeyValueList []KeyValue `json:"keyValueList"` // output
}

func NewImage(ns string, name string, conf string) *Image {
	return &Image{
		Model:        Model{Name: name, namespace: ns},
		Config:       conf,
		CspImageName: "Ubuntu, 18.04",
		Description:  "Canonical, Ubuntu, 18.04 LTS, amd64 bionic",
		GuestOS:      "ubuntu",
		KeyValueList: []KeyValue{},
	}
}

func (self *Image) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/resources/image/%s", self.namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), &self)
}

func (self *Image) POST() error {

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/resources/image?action=registerWithInfo", self.namespace), self, &self)
	if err != nil {
		return err
	}
	return nil
}

func (self *Image) DELETE(ns string) error {

	exist, err := self.GET()
	if err != nil {
		return err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/resources/image/%s", self.namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), model.Status{})
		if err != nil {
			return err
		}
	} else {
		logger.Infof("delete image skip (name=%s, cause=not found)", self.Name)
	}
	return nil
}
