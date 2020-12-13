package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	"github.com/go-resty/resty/v2"
)

type Firewall struct {
	Model
	Config               string          `json:"connectionName"`
	VPCId                string          `json:"vNetId"`
	Description          string          `json:"description"`
	FirewallRules        []FirewallRules `json:"firewallRules"`
	CspSecurityGroupId   string          `json:"cspSecurityGroupId"`   // output
	CspSecurityGroupName string          `json:"cspSecurityGroupName"` // output
	KeyValueList         []KeyValue      `json:"keyValueList"`         // output
}

type FirewallRules struct {
	From      string `json:"fromPort"`
	To        string `json:"toPort"`
	Protocol  string `json:"ipProtocol"`
	Direction string `json:"direction"`
}

func NewFirewall(ns string, name string, conf string) *Firewall {
	return &Firewall{
		Model:  Model{Name: name, namespace: ns},
		Config: conf,
		FirewallRules: []FirewallRules{
			FirewallRules{Protocol: "tcp", Direction: "inbound", From: "1", To: "65535"},
			FirewallRules{Protocol: "udp", Direction: "inbound", From: "1", To: "65535"},
			FirewallRules{Protocol: "icmp", Direction: "inbound", From: "-1", To: "-1"},
		},
	}
}

func (fw *Firewall) GET() (bool, error) {
	// validation
	if err := fw.validate(validation.Validation{}); err != nil {
		return false, err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetBody(fmt.Sprintf(`{"connectionName" : "%s"}`, fw.Config)).
		SetResult(&fw).
		Get(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/securityGroup/%s", fw.namespace, fw.Name))

	if err = fw.response(resp, err); err != nil {
		return false, err
	}
	if resp.StatusCode() == http.StatusNotFound {
		return false, nil
	}

	return true, nil
}

func (fw *Firewall) POST() error {
	// validation
	if err := fw.validate(validation.Validation{}); err != nil {
		return err
	}

	conf := config.Config
	resp, err := resty.New().R().
		SetBasicAuth(conf.Username, conf.Password).
		SetBody(fw).
		SetResult(&fw).
		Post(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/securityGroup", fw.namespace))

	if err = fw.response(resp, err); err != nil {
		return err
	}

	return nil
}

func (fw *Firewall) DELETE(ns string) error {
	// validation
	if err := fw.validate(validation.Validation{}); err != nil {
		return err
	}

	exist, err := fw.GET()
	if err != nil {
		return err
	}
	if exist {
		conf := config.Config
		resp, err := resty.New().R().
			SetBasicAuth(conf.Username, conf.Password).
			SetBody(fmt.Sprintf(`{"connectionName" : "%s"}`, fw.Config)).
			SetResult(TumblebugResult{}).
			Delete(conf.TumblebugUrl + fmt.Sprintf("/ns/%s/resources/securityGroup/%s", fw.namespace, fw.Name))

		if err = fw.response(resp, err); err != nil {
			return err
		}
	} else {
		fmt.Println(fmt.Sprintf("delete firewall skip (name=%s, cause=not found)", fw.Name))
	}

	return nil
}
