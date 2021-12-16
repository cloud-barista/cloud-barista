package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloud-barista/cb-mcks/src/core/model"
	"github.com/cloud-barista/cb-mcks/src/core/model/tumblebug"
	"github.com/cloud-barista/cb-mcks/src/utils/config"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"

	logger "github.com/sirupsen/logrus"
)

type NodeConfigInfo struct {
	model.NodeConfig
	Csp     config.CSP `json:"csp"`
	Role    string     `json:"role"`
	ImageId string     `json:"imageId"`
}

func SetNodeConfigInfos(nodeConfigs []model.NodeConfig, role string) ([]NodeConfigInfo, error) {
	var nodeConfigInfos []NodeConfigInfo

	for _, nodeConfig := range nodeConfigs {
		if nodeConfig.Count < 1 {
			logger.Errorf("node count must be at least one (role=%s, count=%d)", role, nodeConfig.Count)
			return nil, errors.New(fmt.Sprintf("Node count must be at least one. (role=%s, count=%d)", role, nodeConfig.Count))
		}

		conn := tumblebug.NewConnection(nodeConfig.Connection)
		if exists, err := conn.GET(); err != nil {
			return nil, errors.New(fmt.Sprintf("Connection connect error. (role=%s, connection=%s)", role, nodeConfig.Connection))
		} else if !exists {
			logger.Errorf("connection does not exist (role=%s, connection=%s)", role, nodeConfig.Connection)
			return nil, errors.New(fmt.Sprintf("Connection does not exist. (role=%s, connection=%s)", role, nodeConfig.Connection))
		}
		csp, err := GetCSPName(conn.ProviderName)
		if err != nil {
			return nil, err
		}

		region := tumblebug.NewRegion(conn.RegionName)
		if exists, err := region.GET(); err != nil {
			return nil, errors.New(fmt.Sprintf("Region connect error. (role=%s, connection=%s)", role, nodeConfig.Connection))
		} else if !exists {
			logger.Errorf("region does not exist (role=%s, connection=%s, region=%s)", role, nodeConfig.Connection, conn.RegionName)
			return nil, errors.New(fmt.Sprintf("Region does not exist. (region=%s, connection=%s)", role, nodeConfig.Connection))
		}

		imageId, err := GetVmImageId(csp, nodeConfig.Connection, region)
		if err != nil {
			return nil, err
		}

		if err = CheckSpec(csp, nodeConfig.Connection, nodeConfig.Spec, role); err != nil {
			return nil, err
		}

		var nodeConfigInfo NodeConfigInfo
		nodeConfigInfo.Connection = nodeConfig.Connection
		nodeConfigInfo.Count = nodeConfig.Count
		nodeConfigInfo.Spec = nodeConfig.Spec
		nodeConfigInfo.Csp = csp
		nodeConfigInfo.Role = role
		nodeConfigInfo.ImageId = imageId

		nodeConfigInfos = append(nodeConfigInfos, nodeConfigInfo)
	}

	return nodeConfigInfos, nil
}

func GetControlPlaneIPs(VMs []model.VM) []string {
	var IPs []string
	for _, vm := range VMs {
		if vm.Role == config.CONTROL_PLANE {
			IPs = append(IPs, vm.PrivateIP)
		}
	}
	return IPs
}

func GetVmImageName(name string) string {
	tmp := lang.GetOnlyLettersAndNumbers(name)

	return strings.ToLower(tmp)
}

func CheckNamespace(namespace string) error {
	ns := tumblebug.NewNS(namespace)
	if exists, err := ns.GET(); err != nil {
		return err
	} else if !exists {
		logger.Errorf("namespace does not exist (namespace=%s)", namespace)
		return errors.New(fmt.Sprintf("Namespace does not exist. (namespace=%s)", namespace))
	}
	return nil
}

func CheckMcis(namespace string, mcisName string) error {
	mcis := tumblebug.NewMCIS(namespace, mcisName)
	if exists, err := mcis.GET(); err != nil {
		return err
	} else if !exists {
		logger.Errorf("MCIS does not exist (namespace=%s, mcis=%s)", namespace, mcisName)
		return errors.New(fmt.Sprintf("MCIS does not exist. (mcis=%s)", mcisName))
	}
	return nil
}

func CheckClusterStatus(namespace string, clusterName string) error {
	cluster := model.NewCluster(namespace, clusterName)
	if exists, err := cluster.Select(); err != nil {
		return err
	} else if exists == false {
		logger.Errorf("cluster not found(namespace=%s, mcis=%s)", namespace, clusterName)
		return errors.New(fmt.Sprintf("Cluster not found. (namespace=%s, cluster=%s)", namespace, clusterName))
	} else if cluster.Status.Phase != model.ClusterPhaseProvisioned {
		logger.Errorf("cannot add node. cluster status is %s (namespace=%s, mcis=%s)", cluster.Status.Phase, namespace, clusterName)
		return errors.New(fmt.Sprintf("cannot add node. status is '%s'.", cluster.Status.Phase))
	}
	return nil
}

func CheckSpec(csp config.CSP, configName string, specName string, role string) error {
	lookupSpec := tumblebug.NewLookupSpec(configName, specName)
	if err := lookupSpec.LookupSpec(); err != nil {
		return errors.New(fmt.Sprintf("Failed to lookup spec. (csp='%s', spec='%s', cause=%v)", csp, specName, err))
	}

	if lookupSpec.SpiderSpecInfo.Name == "" {
		logger.Errorf("spec '%s' not found (csp=%s)", csp, configName)
		return errors.New(fmt.Sprintf("Failed to find spec. (csp='%s', spec='%s')", csp, specName))
	}

	if role == config.CONTROL_PLANE {
		vCpuCount, err := strconv.Atoi(lookupSpec.SpiderSpecInfo.VCpu.Count)
		if err != nil {
			logger.Errorf("failed to convert vCpu count (csp=%s, spec=%s, cpu=%s)", csp, configName, specName, lookupSpec.SpiderSpecInfo.VCpu.Count)
			return errors.New(fmt.Sprintf("Failed to convert vCpu count. (csp='%s', spec='%s', vCpu.Count=%s)", csp, specName, lookupSpec.SpiderSpecInfo.VCpu.Count))
		}
		if vCpuCount < 2 {
			logger.Errorf("kubernetes control plane node needs 2 vCPU at least (csp=%s, spec=%s, cpu=%d)", csp, configName, specName, vCpuCount)
			return errors.New(fmt.Sprintf("Kubernetes control plane node needs 2 vCPU at least. (csp='%s', spec='%s', cpu=%d)", csp, specName, vCpuCount))
		}
	}

	mem, err := strconv.Atoi(lookupSpec.SpiderSpecInfo.Mem)
	if err != nil {
		logger.Errorf("Failed to convert memory (csp=%s, spec=%s, memory=%s)", csp, specName, lookupSpec.SpiderSpecInfo.Mem)
		return errors.New(fmt.Sprintf("Failed to convert memory. (csp='%s', spec='%s', memory=%s)", csp, specName, lookupSpec.SpiderSpecInfo.Mem))
	}

	gbMem := mem / 1024
	if gbMem < 2 {
		logger.Errorf("kubernetes node needs 2 GiB or more of RAM (csp=%s, spec=%s, memory=%dGB)", csp, specName, gbMem)
		return errors.New(fmt.Sprintf("kubernetes node needs 2 GiB or more of RAM. (csp='%s', spec='%s', memory=%dGB)", csp, specName, gbMem))
	}

	return nil
}
