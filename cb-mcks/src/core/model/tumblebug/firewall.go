package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/cloud-barista/cb-mcks/src/utils/config"
	"github.com/cloud-barista/cb-mcks/src/core/model"
	logger "github.com/sirupsen/logrus"
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

func NewFirewall(csp config.CSP, ns string, name string, conf string) *Firewall {

	fw := &Firewall{
		Model:  Model{Name: name, namespace: ns},
		Config: conf,
		FirewallRules: []FirewallRules{
			{Protocol: "tcp", Direction: "inbound", From: "1", To: "65535"},
			{Protocol: "udp", Direction: "inbound", From: "1", To: "65535"},
		},
	}
	if csp == config.CSP_TENCENT {
		fw.FirewallRules = append(fw.FirewallRules,
			FirewallRules{Protocol: "icmp", Direction: "inbound", From: "ALL"},
			FirewallRules{Protocol: "ALL", Direction: "outbound", From: "ALL"},
		)
	} else {
		fw.FirewallRules = append(fw.FirewallRules, FirewallRules{Protocol: "icmp", Direction: "inbound", From: "-1", To: "-1"})
	}

	return fw
}

func (self *Firewall) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/resources/securityGroup/%s", self.namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), &self)

}

func (self *Firewall) POST() error {

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/resources/securityGroup", self.namespace), self, &self)
	if err != nil {
		return err
	}

	return nil
}

func (self *Firewall) DELETE(ns string) error {

	exist, err := self.GET()
	if err != nil {
		return err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/resources/securityGroup/%s", self.namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), model.Status{})
		if err != nil {
			return err
		}
	} else {
		logger.Infof("delete firewall skip (name=%s, cause=not found)", self.Name)
	}

	return nil
}
