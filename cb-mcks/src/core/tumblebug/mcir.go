package tumblebug

import (
	"fmt"
	"net/http"

	"github.com/beego/beego/v2/core/validation"
	"github.com/cloud-barista/cb-mcks/src/core/app"
)

/* new instance of VPC */
func NewVPC(ns string, name string, conf string, cidrBlock string) *VPC {

	return &VPC{
		Model:     Model{Name: name, Namespace: ns},
		Config:    conf,
		CidrBlock: cidrBlock,
		Subnets: []Subnet{
			{
				Name:      fmt.Sprintf("%s-subnet", conf),
				CidrBlock: cidrBlock},
		},
	}
}

/* new instance of Firewall */
func NewFirewall(csp app.CSP, ns string, name string, conf string) *Firewall {

	fw := &Firewall{
		Model:  Model{Name: name, Namespace: ns},
		Config: conf,
		FirewallRules: []FirewallRules{
			{Protocol: "tcp", Direction: "inbound", From: "1", To: "65535"},
			{Protocol: "udp", Direction: "inbound", From: "1", To: "65535"},
		},
	}
	if csp == app.CSP_TENCENT {
		fw.FirewallRules = append(fw.FirewallRules,
			FirewallRules{Protocol: "icmp", Direction: "inbound", From: "ALL"},
			FirewallRules{Protocol: "ALL", Direction: "outbound", From: "ALL"},
		)
	} else if csp == app.CSP_CLOUDIT {
		fw.FirewallRules = append(fw.FirewallRules, FirewallRules{Protocol: "ALL", Direction: "outbound", From: "ALL"})
	} else {
		fw.FirewallRules = append(fw.FirewallRules, FirewallRules{Protocol: "icmp", Direction: "inbound", From: "-1", To: "-1"})
	}

	return fw
}

/* new instance of VM-Image */
func NewImage(ns string, name string, conf string) *Image {
	return &Image{
		Model:        Model{Name: name, Namespace: ns},
		Config:       conf,
		CspImageName: "Ubuntu, 18.04",
		Description:  "Canonical, Ubuntu, 18.04 LTS, amd64 bionic",
		GuestOS:      "ubuntu",
		KeyValueList: []KeyValue{},
	}
}

/* new instance of VM-Image lookup */
func NewLookupImages(conf string) *LookupImages {
	return &LookupImages{
		Model:          Model{Name: "conf"},
		ConnectionName: conf,
	}
}

/* new instance of VM-Spec. */
func NewSpec(ns string, name string, conf string) *Spec {
	return &Spec{
		Model:  Model{Name: name, Namespace: ns},
		Config: conf,
	}
}

/* new instance of VM-Spec-lookup */
func NewLookupSpec(conf string, spec string) *LookupSpec {
	return &LookupSpec{
		Model:  Model{Name: spec},
		Config: conf,
		Spec:   spec,
	}
}

/* new instance of VM-Specs-lookup */
func NewLookupSpecs(conf string) *LookupSpecs {
	return &LookupSpecs{
		Model:  Model{Name: "conf"},
		Config: conf,
	}
}

/* new instance of SSH-Key */
func NewSSHKey(ns string, name string, conf string) *SSHKey {
	return &SSHKey{
		Model:    Model{Name: name, Namespace: ns},
		Config:   conf,
		Username: VM_USER_ACCOUNT,
	}
}

/* VPC */
func (self *VPC) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/resources/vNet/%s", self.Namespace, self.Name), nil, &self)

}

func (self *VPC) POST() error {

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/resources/vNet", self.Namespace), self, &self)
	if err != nil {
		return err
	}

	return nil

}

func (self *VPC) DELETE() (bool, error) {

	exist, err := self.GET()
	if err != nil {
		return exist, err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/resources/vNet/%s", self.Namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), app.Status{})
		if err != nil {
			return exist, err
		}
	}

	return exist, nil
}

/* Firewall */
func (self *Firewall) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/resources/securityGroup/%s", self.Namespace, self.Name), nil, &self)

}

func (self *Firewall) POST() error {

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/resources/securityGroup", self.Namespace), self, &self)
	if err != nil {
		return err
	}

	return nil
}

func (self *Firewall) DELETE(ns string) (bool, error) {

	exist, err := self.GET()
	if err != nil {
		return exist, err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/resources/securityGroup/%s", self.Namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), app.Status{})
		if err != nil {
			return exist, err
		}
	}

	return exist, nil
}

/* VM-Image */
func (self *Image) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/resources/image/%s", self.Namespace, self.Name), nil, &self)
}

func (self *Image) POST() error {

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/resources/image?action=registerWithInfo", self.Namespace), self, &self)
	if err != nil {
		return err
	}
	return nil
}

func (self *Image) DELETE(ns string) (bool, error) {

	exist, err := self.GET()
	if err != nil {
		return exist, err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/resources/image/%s", self.Namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), app.Status{})
		if err != nil {
			return exist, err
		}
	}
	return exist, nil
}

/* VM-Spec */
func (self *Spec) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/resources/spec/%s", self.Namespace, self.Name), nil, &self)

}

func (self *Spec) POST() error {

	// validation
	valid := validation.Validation{}
	valid.Required(self.CspSpecName, "cspSpecName")
	if err := self.validate(valid); err != nil {
		return err
	}

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/resources/spec", self.Namespace), self, &self)
	if err != nil {
		return err
	}

	return nil

}

func (self *Spec) DELETE(ns string) (bool, error) {

	exist, err := self.GET()
	if err != nil {
		return exist, err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/resources/spec/%s", self.Namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), app.Status{})
		if err != nil {
			return exist, err
		}
	}

	return exist, nil
}

/* SSH-Key */
func (self *SSHKey) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/ns/%s/resources/sshKey/%s", self.Namespace, self.Name), nil, &self)

}

func (self *SSHKey) POST() error {
	// validation
	valid := validation.Validation{}
	valid.Required(self.Username, "username")
	if err := self.validate(valid); err != nil {
		return err
	}

	_, err := self.execute(http.MethodPost, fmt.Sprintf("/ns/%s/resources/sshKey", self.Namespace), self, &self)
	if err != nil {
		return err
	}

	return nil

}

func (self *SSHKey) DELETE(ns string) (bool, error) {

	exist, err := self.GET()
	if err != nil {
		return exist, err
	}
	if exist {
		_, err := self.execute(http.MethodDelete, fmt.Sprintf("/ns/%s/resources/sshKey/%s", self.Namespace, self.Name), fmt.Sprintf(`{"connectionName" : "%s"}`, self.Config), app.Status{})
		if err != nil {
			return exist, err
		}
	}

	return exist, nil
}

/* Look up (VM-Image, VM-Spec) */
func (self *LookupImages) GET() (bool, error) {

	return self.execute(http.MethodPost, "/lookupImages", self, &self)

}

func (spec *LookupSpec) GET() (bool, error) {

	return spec.execute(http.MethodPost, "/lookupSpec", spec, &spec)

}

func (spec *LookupSpecs) GET() (bool, error) {

	return spec.execute(http.MethodPost, "/lookupSpecs", spec, &spec)

}
