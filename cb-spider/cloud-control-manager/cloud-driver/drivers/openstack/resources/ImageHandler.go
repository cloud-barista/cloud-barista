package resources

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	imgsvc "github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"

	call "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/call-log"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
)

const (
	Image = "IMAGE"
)

type OpenStackImageHandler struct {
	Client      *gophercloud.ServiceClient
	ImageClient *gophercloud.ServiceClient
}

func setterImage(image images.Image) *irs.ImageInfo {
	imageInfo := &irs.ImageInfo{
		IId: irs.IID{
			NameId:   image.ID,
			SystemId: image.ID,
		},
		GuestOS: image.Name,
		Status: image.Status,
	}

	// 메타 정보 등록
	var metadataList []irs.KeyValue
	for key, val := range image.Metadata {
		metadata := irs.KeyValue{
			Key:   key,
			Value: fmt.Sprintf("%v", val),
		}
		metadataList = append(metadataList, metadata)
	}
	imageInfo.KeyValueList = metadataList

	return imageInfo
}

func (imageHandler *OpenStackImageHandler) CreateImage(imageReqInfo irs.ImageReqInfo) (irs.ImageInfo, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(imageHandler.Client.IdentityEndpoint, call.VMIMAGE, imageReqInfo.IId.NameId, "CreateImage()")

	// @TODO: Image 생성 요청 파라미터 정의 필요
	type ImageReqInfo struct {
		Name            string
		ContainerFormat string
		DiskFormat      string
	}
	reqInfo := ImageReqInfo{
		Name:            imageReqInfo.IId.NameId,
		ContainerFormat: "bare",
		DiskFormat:      "iso",
	}

	createOpts := imgsvc.CreateOpts{
		Name:            reqInfo.Name,
		ContainerFormat: reqInfo.ContainerFormat,
		DiskFormat:      reqInfo.DiskFormat,
	}

	// Check Image file exists
	rootPath := os.Getenv("CBSPIDER_ROOT")
	imageFilePath := fmt.Sprintf("%s/image/%s.iso", rootPath, reqInfo.Name)
	if _, err := os.Stat(imageFilePath); os.IsNotExist(err) {
		createErr := errors.New(fmt.Sprintf("Image files in path %s not exist", imageFilePath))
		cblogger.Error(createErr.Error())
		LoggingError(hiscallInfo, createErr)
		return irs.ImageInfo{}, createErr
	}

	// Create Image
	start := call.Start()
	image, err := imgsvc.Create(imageHandler.ImageClient, createOpts).Extract()
	if err != nil {
		cblogger.Error(err.Error())
		LoggingError(hiscallInfo, err)
		return irs.ImageInfo{}, err
	}
	LoggingInfo(hiscallInfo, start)

	// Upload Image file
	imageBytes, err := ioutil.ReadFile(imageFilePath)
	if err != nil {
		cblogger.Error(err.Error())
		LoggingError(hiscallInfo, err)
		return irs.ImageInfo{}, err
	}
	result := imagedata.Upload(imageHandler.ImageClient, image.ID, bytes.NewReader(imageBytes))
	if result.Err != nil {
		cblogger.Error(result.Err.Error())
		LoggingError(hiscallInfo, err)
		return irs.ImageInfo{}, err
	}

	// 생성된 Imgae 정보 리턴
	mappedImageInfo := images.Image{
		ID:       image.ID,
		Created:  image.CreatedAt.String(),
		MinDisk:  image.MinDiskGigabytes,
		MinRAM:   image.MinRAMMegabytes,
		Name:     image.Name,
		Status:   string(image.Status),
		Updated:  image.UpdatedAt.String(),
		Metadata: image.Properties,
	}
	imageInfo := setterImage(mappedImageInfo)
	return *imageInfo, nil
}

func (imageHandler *OpenStackImageHandler) ListImage() ([]*irs.ImageInfo, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(imageHandler.Client.IdentityEndpoint, call.VMIMAGE, Image, "ListImage()")

	start := call.Start()
	pager, err := images.ListDetail(imageHandler.Client, images.ListOpts{}).AllPages()
	if err != nil {
		cblogger.Error(err.Error())
		LoggingError(hiscallInfo, err)
		return nil, err
	}
	LoggingInfo(hiscallInfo, start)

	imageList, err := images.ExtractImages(pager)
	if err != nil {
		cblogger.Error(err.Error())
		LoggingError(hiscallInfo, err)
		return nil, err
	}

	imageInfoList := make([]*irs.ImageInfo, len(imageList))
	for i, img := range imageList {
		imageInfo := setterImage(img)
		imageInfoList[i] = imageInfo
	}
	return imageInfoList, nil
}

func (imageHandler *OpenStackImageHandler) GetImage(imageIID irs.IID) (irs.ImageInfo, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(imageHandler.Client.IdentityEndpoint, call.VMIMAGE, imageIID.NameId, "GetImage()")

	start := call.Start()
	image, err := imageHandler.getRawImage(imageIID)
	if err != nil {
		cblogger.Error(err.Error())
		LoggingError(hiscallInfo, err)
		return irs.ImageInfo{}, err
	}
	LoggingInfo(hiscallInfo, start)

	imageInfo := setterImage(*image)
	return *imageInfo, nil
}

func (imageHandler *OpenStackImageHandler) DeleteImage(imageIID irs.IID) (bool, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(imageHandler.Client.IdentityEndpoint, call.VMIMAGE, imageIID.NameId, "DeleteImage()")

	image, err := imageHandler.getRawImage(imageIID)
	if err != nil {
		cblogger.Error(err.Error())
		LoggingError(hiscallInfo, err)
		return false, err
	}
	start := call.Start()
	err = images.Delete(imageHandler.Client, image.ID).ExtractErr()
	if err != nil {
		cblogger.Error(err.Error())
		LoggingError(hiscallInfo, err)
		return false, err
	}
	LoggingInfo(hiscallInfo, start)
	return true, nil
}
func (imageHandler *OpenStackImageHandler) getRawImage(imageIId irs.IID) (*images.Image, error) {
	if !CheckIIDValidation(imageIId) {
		return nil, errors.New("invalid IID")
	}
	return images.Get(imageHandler.Client, imageIId.SystemId).Extract()
}

