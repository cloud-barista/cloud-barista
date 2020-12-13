package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	"github.com/go-resty/resty/v2"
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

func (image *Image) GET() (bool, error) {
	// validation
	if err := image.validate(validation.Validation{}); err != nil {
		return false, err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetBody(fmt.Sprintf(`{"connectionName" : "%s"}`, image.Config)).
		SetResult(&image).
		Get(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/image/%s", image.namespace, image.Name))

	if e := image.response(resp, err); e != nil {
		return false, e
	}
	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}

func (image *Image) POST() error {
	// validation
	if err := image.validate(validation.Validation{}); err != nil {
		return err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetBody(image).
		SetResult(&image).
		Post(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/image?action=registerWithInfo", image.namespace))

	if e := image.response(resp, err); e != nil {
		return e
	}

	return nil
}

func (image *Image) DELETE(ns string) error {
	// validation
	if err := image.validate(validation.Validation{}); err != nil {
		return err
	}

	exist, err := image.GET()
	if err != nil {
		return err
	}
	if exist {
		conf := config.Config
		resp, err := resty.New().R().
			SetBasicAuth(conf.Username, conf.Password).
			SetBody(fmt.Sprintf(`{"connectionName" : "%s"}`, image.Config)).
			SetResult(TumblebugResult{}).
			Delete(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/image/%s", image.namespace, image.Name))

		if e := image.response(resp, err); e != nil {
			return e
		}

	} else {
		fmt.Println(fmt.Sprintf("delete image skip (name=%s, cause=not found)", image.Name))
	}
	return nil
}
