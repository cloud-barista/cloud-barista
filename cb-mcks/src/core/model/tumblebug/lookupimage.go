package tumblebug

import (
	"net/http"
)

type LookupImages struct {
	Model
	ConnectionName string            `json:"connectionName"`
	Image          []SpiderImageInfo `json:"image"`
}

type SpiderImageInfo struct {
	// Fields for request
	Name string
	// Fields for response
	IId          IID    // {NameId, SystemId}
	GuestOS      string // Windows7, Ubuntu etc.
	Status       string // available, unavailable
	KeyValueList []KeyValue
}

// Spider 2020-03-30 https://github.com/cloud-barista/cb-spider/blob/master/cloud-control-manager/cloud-driver/interfaces/resources/IId.go
type IID struct {
	NameId   string // NameID by user
	SystemId string // SystemID by CloudOS
}

type LookupImagesInfo struct {
	NameId   string
	FilterId string
}

func NewLookupImages(conf string) *LookupImages {
	return &LookupImages{
		ConnectionName: conf,
	}
}

func (image *LookupImages) LookupImages() error {
	_, err := image.execute(http.MethodPost, "/lookupImages", image, &image)
	if err != nil {
		return err
	}

	return nil
}
