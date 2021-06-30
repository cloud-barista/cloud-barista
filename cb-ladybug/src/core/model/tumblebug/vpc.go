package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-ladybug/src/core/model"
	logger "github.com/sirupsen/logrus"
)

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
			{
				Name:      fmt.Sprintf("%s-subnet", name),
				CidrBlock: "192.168.1.0/24"},
		},
	}
}

func (self *VPC) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/resources/vNet/%s", self.namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), &self)

}

func (self *VPC) POST() error {

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/resources/vNet", self.namespace), self, &self)
	if err != nil {
		return err
	}

	return nil

}

func (self *VPC) DELETE() error {

	exist, err := self.GET()
	if err != nil {
		return err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/resources/vNet/%s", self.namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), model.Status{})
		if err != nil {
			return err
		}
	} else {
		logger.Infof("delete vpc skip.. (name=%s, cause=not found)", self.Name)
	}

	return nil
}
