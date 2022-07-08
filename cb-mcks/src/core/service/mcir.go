package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloud-barista/cb-mcks/src/core/app"
	"github.com/cloud-barista/cb-mcks/src/core/model"
	"github.com/cloud-barista/cb-mcks/src/core/tumblebug"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"

	logger "github.com/sirupsen/logrus"
)

type MCIR struct {
	namespace    string
	csp          app.CSP
	role         app.ROLE
	config       string //prameter
	spec         string //prameter
	vmCount      int    //prameter
	credential   string
	subnetName   string
	vpcName      string
	firewallName string
	sshkeyName   string
	imageName    string
	specName     string
	region       string
	zone         string
}

func NewMCIR(namespace string, role app.ROLE, nodeSetReq app.NodeSetReq) *MCIR {

	specName := strings.ToLower(lang.ReplaceAll(nodeSetReq.Spec, []string{".", "_", " "}, "-"))

	return &MCIR{
		namespace:    namespace,
		role:         role,
		config:       nodeSetReq.Connection,
		spec:         nodeSetReq.Spec,
		vmCount:      nodeSetReq.Count,
		vpcName:      fmt.Sprintf("%s-vpc", nodeSetReq.Connection),
		firewallName: fmt.Sprintf("%s-sg", nodeSetReq.Connection),
		sshkeyName:   fmt.Sprintf("%s-sshkey", nodeSetReq.Connection),
		imageName:    fmt.Sprintf("%s-ubuntu1804", nodeSetReq.Connection),
		specName:     fmt.Sprintf("%s-%s-spec", nodeSetReq.Connection, specName),
	}
}

/* create a MCIR (vpc, firewall, ssk-key, vm-spec, vm-image) if there is not exist */
func (self *MCIR) CreateIfNotExist() (model.ClusterReason, string) {

	// validate a connection info.
	connection := tumblebug.NewConnection(self.config)
	if exists, err := connection.GET(); err != nil {
		return model.InvalidMCIRReason, fmt.Sprintf("Failed to get a connection info. (%s)", self.config)
	} else if !exists {
		return model.InvalidMCIRReason, fmt.Sprintf("Connection does not exist. (%s)", self.config)
	}
	self.csp = app.CSP(strings.ToLower(connection.ProviderName))

	//validate a CSP
	exists := false
	for _, c := range []string{string(app.CSP_AWS), string(app.CSP_GCP), string(app.CSP_AZURE), string(app.CSP_ALIBABA), string(app.CSP_TENCENT), string(app.CSP_OPENSTACK), string(app.CSP_IBM), string(app.CSP_CLOUDIT)} {
		if string(self.csp) == c {
			exists = true
			break
		}
	}
	if exists == false {
		return model.InvalidMCIRReason, fmt.Sprintf("The CSP '%s' is not supported", connection.ProviderName)
	}

	// validation a spec.
	if err := self.verifySpec(); err != nil {
		return model.InvalidMCIRReason, err.Error()
	}

	// get a region
	region := tumblebug.NewRegion(connection.RegionName)
	if exists, err := region.GET(); err != nil {
		return model.InvalidMCIRReason, fmt.Sprintf("Failed to get a region data. (cause='%v')", err)
	} else if !exists {
		return model.InvalidMCIRReason, fmt.Sprintf("Region does not exist (%s)", connection.RegionName)
	}
	for _, r := range region.KeyValueInfoList {
		if r.Key == "Region" {
			self.region = r.Value
		} else if r.Key == "Zone" {
			self.zone = r.Value
		}

	}

	// Create a VPC
	vpc := tumblebug.NewVPC(self.namespace, self.vpcName, self.config, getCSPCidrBlock(self.csp))
	exists, err := vpc.GET()
	if err != nil {
		return model.CreateVpcFailedReason, fmt.Sprintf("Failed to create a VPC. (cause='%v')", err)
	}
	if exists {
		logger.Infof("[%s] VPC has been reused. (%s)", self.config, self.vpcName)
	} else {
		if err = vpc.POST(); err != nil {
			return model.CreateVpcFailedReason, fmt.Sprintf("Failed to create a VPC. (cause='%v')", err)
		}
		logger.Infof("[%s] VPC creation has been completed. (%s)", self.config, self.vpcName)
	}
	self.subnetName = vpc.Subnets[0].Name

	// Create a Firewall
	fw := tumblebug.NewFirewall(self.csp, self.namespace, self.firewallName, self.config)
	fw.VPCId = self.vpcName
	exists, err = fw.GET()
	if err != nil {
		return model.CreateSecurityGroupFailedReason, fmt.Sprintf("Failed to create a Firewall Rules. (cause='%v')", err)
	}
	if exists {
		logger.Infof("[%s] Firewall has been reused. (%s)", self.config, self.firewallName)
	} else {
		if err = fw.POST(); err != nil {
			return model.CreateSecurityGroupFailedReason, fmt.Sprintf("Failed to create a Firewall Rules. (cause='%v')", err)
		}
		logger.Infof("[%s] Firewall creation has been completed. (%s)", self.config, self.firewallName)
	}

	// Create a SSH-Key
	sshKey := tumblebug.NewSSHKey(self.namespace, self.sshkeyName, self.config)
	exists, err = sshKey.GET()
	if err != nil {
		return model.CreateSSHKeyFailedReason, fmt.Sprintf("Failed to create a SSH-Key. (cause='%v')", err)
	}
	if exists {
		logger.Infof("[%s] SSH-Key has been reused. (%s)", self.config, self.sshkeyName)
	} else {
		if err = sshKey.POST(); err != nil {
			return model.CreateSSHKeyFailedReason, fmt.Sprintf("Failed to create a SSH-Key. (cause='%v')", err)
		}
		logger.Infof("[%s] SSH-Key creation has been completed. (%s)", self.config, self.sshkeyName)
	}
	self.credential = sshKey.PrivateKey

	// Create a Image
	imageId, err := getCSPImageId(self.csp, self.config, region)
	if err != nil {
		return model.InvalidMCIRReason, err.Error()
	}
	image := tumblebug.NewImage(self.namespace, self.imageName, self.config)
	image.CspImageId = imageId
	exists, err = image.GET()
	if err != nil {
		return model.CreateVmImageFailedReason, fmt.Sprintf("Failed to create a Image. (cause='%v')", err)
	}
	if exists {
		logger.Infof("[%s] VM-Image has been reused. (%s)", self.config, self.imageName)
	} else {
		if err = image.POST(); err != nil {
			return model.CreateVmImageFailedReason, fmt.Sprintf("Failed to create a Image. (cause='%v')", err)
		}
		logger.Infof("[%s] VM-Image creation has been completed. (%s)", self.config, self.imageName)
	}

	// Create a spec.
	spec := tumblebug.NewSpec(self.namespace, self.specName, self.config)
	spec.CspSpecName = self.spec
	exists, err = spec.GET()
	if err != nil {
		return model.CreateVmSpecFailedReason, fmt.Sprintf("Failed to create a VM Spec. (cause='%v')", err)
	}
	if exists {
		logger.Infof("[%s] VM-Spec has been reused. (%s)", self.config, self.spec)
	} else {
		if err = spec.POST(); err != nil {
			return model.CreateVmSpecFailedReason, fmt.Sprintf("Failed to create a VM Spec. (cause='%v')", err)
		}
		logger.Infof("[%s] VM-Spec creation has been completed. (%s)", self.config, self.spec)
	}

	return "", ""
}

/* new a VM template */
func (self *MCIR) NewVM(namespace string, name string, mcisName string) tumblebug.VM {
	vm := tumblebug.NewVM(namespace, name, mcisName)
	vm.Config = self.config
	vm.VPC = self.vpcName
	vm.Subnet = self.subnetName
	vm.Firewalls = []string{self.firewallName}
	vm.SSHKey = self.sshkeyName
	vm.Image = self.imageName
	vm.Spec = self.specName
	return *vm
}

/* verify - cpus & momories & look-up(exists) */
func (self *MCIR) verifySpec() error {

	lookupSpec := tumblebug.NewLookupSpec(self.config, self.spec)
	if exist, err := lookupSpec.GET(); err != nil {
		return errors.New(fmt.Sprintf("Failed to lookup spec. (csp=%s, spec=%s, cause='%v')", self.csp, self.spec, err))
	} else if !exist {
		return errors.New(fmt.Sprintf("Could not be found a spec '%s'. (connection=%s, csp=%s)", self.spec, self.config, self.csp))
	}

	if self.role == app.CONTROL_PLANE {
		vCpuCount, err := strconv.Atoi(lookupSpec.CPU.Count)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to convert cpu count. (csp=%s, spec=%s, cpu=%s)", self.csp, self.spec, lookupSpec.CPU.Count))
		}
		if vCpuCount < 2 {
			return errors.New(fmt.Sprintf("Kubernetes control plane node needs 2 cpu at least. (csp=%s, spec=%s, cpu=%d)", self.csp, self.spec, vCpuCount))
		}
	}

	mem, err := strconv.Atoi(lookupSpec.Memory)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to convert memory. (csp=%s, spec=%s, memory=%s)", self.csp, self.spec, lookupSpec.Memory))
	}

	gbMem := mem / 1024
	if gbMem < 2 {
		return errors.New(fmt.Sprintf("kubernetes node needs 2 GiB or more of RAM. (csp=%s, spec=%s, memory=%dGB)", self.csp, self.spec, gbMem))
	}

	return nil
}

func VerifySpecList(configName string, controlPlane string, cpumin int, cpumax int, memorymin int, memorymax int) (SpecList, error) {
	if controlPlane == "Y" {
		if cpumin < 2 {
			return SpecList{}, errors.New(fmt.Sprintf("Kubernetes control plane node needs 2 cpu at least. (cpu-min=%d)", cpumin))
		}
		if memorymin < 2 {
			return SpecList{}, errors.New(fmt.Sprintf("Kubernetes control plane node needs 2 memory at least. (memory-min=%d)", memorymin))
		}
	}

	if cpumin > cpumax {
		return SpecList{}, errors.New(fmt.Sprintf("The cpu-max must be greater than or equal to the cpu-min. (cpu-min=%d cpu-max=%d)", cpumin, cpumax))
	}

	if memorymin > memorymax {
		return SpecList{}, errors.New(fmt.Sprintf("The memory-max must be greater than or equal to the memory-min. (memory-min=%d memory-max=%d)", memorymin, memorymax))
	}

	lookupSpecs := tumblebug.NewLookupSpecs(configName)

	if exist, err := lookupSpecs.GET(); err != nil {
		return SpecList{}, errors.New(fmt.Sprintf("Failed to lookup spec. (connection=%s, cause='%v')", configName, err))
	} else if !exist {
		return SpecList{}, errors.New(fmt.Sprintf("Could not be found a specList. (connection=%s)", configName))
	} else {
		var filterSpec []Vmspecs

		for _, Vmspec := range lookupSpecs.Vmspecs {
			cpuCount, err := strconv.Atoi(Vmspec.CPU.Count)
			if err != nil {
				return SpecList{}, errors.New(fmt.Sprintf("Failed to convert CPU. (CPU=%s)", Vmspec.CPU.Count))
			}
			mem, err := strconv.Atoi(Vmspec.Memory)
			if err != nil {
				return SpecList{}, errors.New(fmt.Sprintf("Failed to convert memory. (memory=%s)", Vmspec.Memory))
			}
			mem /= 1024
			if cpuCount >= cpumin && cpumax >= cpuCount && mem >= memorymin && memorymax >= mem {
				VMSpecList := Vmspecs{Name: Vmspec.Name, Memory: strconv.Itoa(mem), CPU: Vmspec.CPU}
				filterSpec = append(filterSpec, VMSpecList)
			}
		}

		SpecList := SpecList{Kind: "specList", Config: configName, Vmspecs: filterSpec}

		return SpecList, nil

	}
}
