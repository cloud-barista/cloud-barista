package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/cloud-barista/cb-mcks/src/core/app"
	"github.com/cloud-barista/cb-mcks/src/core/tumblebug"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"
)

const (
	GCP_IMAGE_ID     = "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-1804-bionic-v20201014"
	AZURE_IMAGE_ID   = "Canonical:UbuntuServer:18.04-LTS:latest"
	ALIBABA_IMAGE_ID = "ubuntu_18_04_x64_20G_alibase_20210521.vhd"
	TENCENT_IMAGE_ID = "img-pi0ii46r"
	CLOUDIT_IMAGE_ID = "Ubuntu 18.04"
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

var ibmImageMap = map[string]string{
	"us-south": "r006-9de77234-3189-42f8-982d-f2266477cfe0", //미국 남부
	"br-sao":   "r042-92d1cd12-f014-4b9a-abf8-c5ca6494a9e5", //브라질
	"us-east":  "r014-dc446598-a1b5-41c3-a1d6-add3afaf264e", //미국 동부
	"eu-de":    "r010-1f68eb2d-f35c-4959-8f4b-2b2f9cf78102", //독일
	"ca-tor":   "r038-e92647cf-8be9-438a-b94c-251cc86bc99a", //캐나다
	"eu-gb":    "r018-1d7417c6-893e-49d4-b14d-9643d6b29812", //영국
	"au-syd":   "r026-a8c25ce6-0ca1-43e9-9b41-411c6217b8b8", //호주
	"jp-osa":   "r034-522c639c-52e1-4cab-8dfb-bc0fb9f6f577", //일본 (오사카)
	"jp-tok":   "r022-61fdadec-6b03-4bd2-bfca-62cd16f5673f", //일본 (도쿄)
}

// get a cidr-block
func getCSPCidrBlock(csp app.CSP) string {

	switch csp {
	case app.CSP_AWS:
		return fmt.Sprintf("192.168.%d.0/24", 10+rand.Intn(10))
	case app.CSP_GCP:
		return fmt.Sprintf("192.168.%d.0/24", 20+rand.Intn(10))
	case app.CSP_AZURE:
		return fmt.Sprintf("192.168.%d.0/24", 30+rand.Intn(10))
	case app.CSP_ALIBABA:
		return fmt.Sprintf("192.168.%d.0/24", 40+rand.Intn(10))
	case app.CSP_TENCENT:
		return fmt.Sprintf("192.168.%d.0/24", 50+rand.Intn(10))
	case app.CSP_IBM:
		return fmt.Sprintf("192.168.%d.0/24", 60+rand.Intn(10))
	case app.CSP_OPENSTACK:
		return fmt.Sprintf("192.168.%d.0/24", 70+rand.Intn(10))
	case app.CSP_CLOUDIT:
		return "10.0.244.0/22"
	}

	return "192.168.255.0/24"
}

// get a vm-image-id
func getCSPImageId(csp app.CSP, configName string, region *tumblebug.Region) (string, error) {

	if csp == app.CSP_GCP {
		return GCP_IMAGE_ID, nil
	} else if csp == app.CSP_AZURE {
		return AZURE_IMAGE_ID, nil
	} else if csp == app.CSP_ALIBABA {
		return ALIBABA_IMAGE_ID, nil
	} else if csp == app.CSP_TENCENT {
		return TENCENT_IMAGE_ID, nil
	} else if csp == app.CSP_CLOUDIT {
		return CLOUDIT_IMAGE_ID, nil
	} else if csp == app.CSP_OPENSTACK {
		// openstack : lookupImages를 통해 사용자가 등록한 이미지를 검색하여, 이미지 이름에 'ubuntu'와 '1804'가 포함된 이미지 정보 가져오기
		lookupImages := tumblebug.NewLookupImages(configName)
		if exist, err := lookupImages.GET(); err != nil {
			return "", errors.New(fmt.Sprintf("Failed to lookup a images on openstack. (connection=%s, cause='%v')", configName, err))
		} else if !exist {
			return "", errors.New(fmt.Sprintf("Could not be found a image on openstack. (connection=%s)", configName))
		}

		for _, image := range lookupImages.Images {
			id := strings.ToLower(lang.GetOnlyLettersAndNumbers(image.IId.NameId))
			if strings.Contains(id, "ubuntu") && strings.Contains(id, "1804") {
				return image.IId.NameId, nil
			}
		}

		return "", errors.New(fmt.Sprintf("Could not be found a ubuntu 18.04 image on openstack. please create an image based on Ubuntu 18.04 The image name must include 'ubuntu' and '18.04'. (connection=%s)", configName))

	} else if csp == app.CSP_AWS {
		// AWS : 리전별 AMI 가져오기
		regionName := ""
		for _, info := range region.KeyValueInfoList {
			if info.Key == "Region" {
				regionName = info.Value //get region name
				break
			}
		}

		if regionName == "" {
			return "", errors.New(fmt.Sprintf("Could not be found a AMI on AWS. (cause = region name is empty, connection=%s, region=%s)", configName, region.RegionName))
		}

		// TODO [update/hard-coding] region별 image id
		imageId := imageMap[regionName]
		if imageId == "" {
			return "", errors.New(fmt.Sprintf("Could not be found a AMI on AWS image map. (connection=%s, region=%s)", configName, regionName))
		}

		return imageId, nil

	} else if csp == app.CSP_IBM {
		// IBM : 리전별 image 가져오기
		regionName := ""
		for _, info := range region.KeyValueInfoList {
			if info.Key == "Region" {
				regionName = info.Value //get region name
				break
			}
		}

		if regionName == "" {
			return "", errors.New(fmt.Sprintf("Could not be found a image on IBM. (cause = region name is empty, connection=%s, region=%s)", configName, region.RegionName))
		}

		// TODO [update/hard-coding] region별 image id
		imageId := ibmImageMap[regionName]
		if imageId == "" {
			return "", errors.New(fmt.Sprintf("Could not be found a image on IBM image map. (connection=%s, region=%s)", configName, regionName))
		}

		return imageId, nil
	} else {
		return "", errors.New(fmt.Sprintf("CSP '%s' is not supported. (could not be found 'vm-machine-image')", csp))
	}

}
