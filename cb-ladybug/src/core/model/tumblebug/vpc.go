package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	"github.com/go-resty/resty/v2"
)

type TumblebugResult struct {
	Message string `json:"message"`
}

type VPC struct {
	Model
	Config       string     `json:"connectionName"`
	CidrBlock    string     `json:"cidrBlock"`
	Subnets      []Subnet   `json:"subnetInfoList"`
	Description  string     `json:"description"`
	CspVNetId    string     `json:"cspVNetId"`    // output
	CspVNetName  string     `json:"cspVNetName"`  // output
	Status       string     `json:"status"`       // output
	KeyValueList []KeyValue `json:"keyValueList"` // output
}

type Subnet struct {
	Name      string `json:"Name"`
	CidrBlock string `json:"IPv4_CIDR"`
}

func NewVPC(ns string, name string, conf string) *VPC {
	return &VPC{
		Model:     Model{Name: name, namespace: ns},
		Config:    conf,
		CidrBlock: "192.168.0.0/16",
		Subnets: []Subnet{
			Subnet{
				Name:      fmt.Sprintf("%s-subnet", name),
				CidrBlock: "192.168.1.0/24"},
		},
	}
}

func (vpc *VPC) GET() (bool, error) {
	// validation
	if err := vpc.validate(validation.Validation{}); err != nil {
		return false, err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetBody(fmt.Sprintf(`{"connectionName" : "%s"}`, vpc.Config)).
		SetResult(&vpc).
		Get(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/vNet/%s", vpc.namespace, vpc.Name))

	if err = vpc.response(resp, err); err != nil {
		return false, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}

func (vpc *VPC) POST() error {
	// validation
	if err := vpc.validate(validation.Validation{}); err != nil {
		return err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetBody(vpc).
		SetResult(&vpc).
		Post(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/vNet", vpc.namespace))

	if err = vpc.response(resp, err); err != nil {
		return err
	}

	return nil
}

func (vpc *VPC) DELETE() error {
	// validation
	if err := vpc.validate(validation.Validation{}); err != nil {
		return err
	}

	exist, err := vpc.GET()
	if err != nil {
		return err
	}
	if exist {
		conf := config.Config
		resp, err := resty.New().R().
			SetBasicAuth(conf.Username, conf.Password).
			SetBody(fmt.Sprintf(`{"connectionName" : "%s"}`, vpc.Config)).
			SetResult(TumblebugResult{}).
			Delete(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/vNet/%s", vpc.namespace, vpc.Name))

		if err = vpc.response(resp, err); err != nil {
			return err
		}
	} else {
		fmt.Println(fmt.Sprintf("delete vpc skip.. (name=%s, cause=not found)", vpc.Name))
	}

	return nil
}
