package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cloud-barista/cb-mcks/src/core/model/tumblebug"
	"github.com/cloud-barista/cb-mcks/src/utils/config"
	logger "github.com/sirupsen/logrus"
)

const (
	GCP_IMAGE_ID     = "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-1804-bionic-v20201014"
	AZURE_IMAGE_ID   = "Canonical:UbuntuServer:18.04-LTS:latest"
	ALIBABA_IMAGE_ID = "ubuntu_18_04_x64_20G_alibase_20210521.vhd"
	TENCENT_IMAGE_ID = "img-pi0ii46r"
)

// region별 AMI :  (AMI 이름 : ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-20200908, 소유자:099720109477 )
var imageMap = map[string]string{
	"us-east-1":      "ami-0817d428a6fb68645", //미국 동부 (버지니아 북부)
	"us-east-2":      "ami-0e82959d4ed12de3f", //미국 동부 (오하이오)
	"us-west-1":      "ami-03fac5402e10ea93b", //미국서부 (캘리포니아)
	"us-west-2":      "ami-07a29e5e945228fa1", //미국서부 (오래곤)
	"af-south-1":     "",                      //아프리카 (케이프타운)
	"ap-east-1":      "",                      //아시아 태평양 (홍콩)
	"ap-south-1":     "ami-03f0fd1a2ba530e75", //아시아 태평양 (뭄바이)
	"ap-northeast-2": "ami-064ab8637cf33f1bb", //아시아 태평양 (서울)
	"ap-southeast-1": "ami-0c8e97a27be37adfd", //아시아 태평양 (싱가포르)
	"ap-southeast-2": "ami-099c1869f33464fde", //아시아 태평양 (시드니)
	"ap-northeast-1": "ami-02b658ac34935766f", //아시아 태평양 (도쿄)
	"ap-northeast-3": "ami-0b627457758376b52", //아시아 태평양 (오사카)
	"ca-central-1":   "ami-0c27a26eca5dc74fc", //캐나다 (중부)
	"eu-central-1":   "ami-092391a11f8aa4b7b", //유럽 (프랑크푸르트)
	"eu-west-1":      "ami-0823c236601fef765", //유럽 (아일랜드)
	"eu-west-2":      "ami-09a1e275e350acf38", //유럽 (런던)
	"eu-south-1":     "",                      //유럽 (밀라노)
	"eu-west-3":      "ami-014d8dccd70fd2632", //유럽 (파리)
	"eu-north-1":     "ami-0ede7f804d699ea83", //유럽 (스톡홀름)
	"me-south-1":     "",                      //중동 (바레인)
	"sa-east-1":      "ami-0fd2c3d373788b726", //남아메리카 (상파울루)
}

// get vm image-id
func GetVmImageId(csp config.CSP, configName string, region *tumblebug.Region) (string, error) {

	if csp == config.CSP_GCP {
		return GCP_IMAGE_ID, nil
	} else if csp == config.CSP_AZURE {
		return AZURE_IMAGE_ID, nil
	} else if csp == config.CSP_ALIBABA {
		return ALIBABA_IMAGE_ID, nil
	} else if csp == config.CSP_TENCENT {
		return TENCENT_IMAGE_ID, nil
	} else if csp == config.CSP_OPENSTACK {
		// openstack : lookupImages를 통해 사용자가 등록한 이미지를 검색하여, 이미지 이름에 'ubuntu'와 '1804'가 포함된 이미지 정보 가져오기
		lookupImages := tumblebug.NewLookupImages(configName)
		err := lookupImages.LookupImages()
		if err != nil {
			return "", errors.New(fmt.Sprintf("failed to lookup images on openstack (connection='%s', cause=%v)", configName, err))
		}

		var filterImages []tumblebug.LookupImagesInfo
		for _, li := range lookupImages.Image {
			imageInfo := tumblebug.LookupImagesInfo{
				NameId:   li.IId.NameId,
				FilterId: GetVmImageName(li.IId.NameId),
			}
			filterImages = append(filterImages, imageInfo)
		}

		for _, fi := range filterImages {
			if strings.Contains(fi.FilterId, "ubuntu") && strings.Contains(fi.FilterId, "1804") {
				return fi.NameId, nil
			}
		}
		return "", errors.New(fmt.Sprintf("request not found ubuntu 18.04 image on openstack. please create an image based on Ubuntu 18.04 The image name must include 'ubuntu' and '18.04'. (connection='%s')", configName))

	} else if csp == config.CSP_AWS {
		// AWS : 리전별 AMI 가져오기
		regionName := ""
		for _, info := range region.KeyValueInfoList {
			if info.Key == "Region" {
				regionName = info.Value //get region name
				break
			}
		}

		if regionName == "" {
			return "", errors.New(fmt.Sprintf("request not found AMI on AWS (cause = region name is empty, connection='%s', region name='%s')", configName, region.RegionName))
		}

		// TODO [update/hard-coding] region별 image id
		imageId := imageMap[regionName]
		if imageId == "" {
			return "", errors.New(fmt.Sprintf("request not found AMI on AWS image map (connection='%s', region='%s')", configName, regionName))
		}

		logger.Infof("AMI find OK (ami='%s', region='%s')", imageId, regionName)
		return imageId, nil

	} else {
		return "", errors.New(fmt.Sprintf("CSP '%s' is not supported (Not found \"vm-machine-image\")", csp))
	}

}

// get CSP Name
func GetCSPName(providerName string) (config.CSP, error) {

	switch strings.ToLower(providerName) {
	case string(config.CSP_AWS):
		return config.CSP_AWS, nil
	case string(config.CSP_GCP):
		return config.CSP_GCP, nil
	case string(config.CSP_AZURE):
		return config.CSP_AZURE, nil
	case string(config.CSP_ALIBABA):
		return config.CSP_ALIBABA, nil
	case string(config.CSP_TENCENT):
		return config.CSP_TENCENT, nil
	case string(config.CSP_OPENSTACK):
		return config.CSP_OPENSTACK, nil
	}
	return "", errors.New(providerName + "is not supported")
}
