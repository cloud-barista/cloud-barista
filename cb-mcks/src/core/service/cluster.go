package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-barista/cb-mcks/src/core/model"
	"github.com/cloud-barista/cb-mcks/src/core/model/tumblebug"
	"github.com/cloud-barista/cb-mcks/src/utils/config"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"
	"golang.org/x/sync/errgroup"

	logger "github.com/sirupsen/logrus"
)

func ListCluster(namespace string) (*model.ClusterList, error) {

	// validate namespace
	if err := CheckNamespace(namespace); err != nil {
		return nil, err
	}

	clusters := model.NewClusterList(namespace)
	if err := clusters.SelectList(); err != nil {
		return nil, err
	}

	return clusters, nil
}

func GetCluster(namespace string, clusterName string) (*model.Cluster, error) {

	// validate namespace
	if err := CheckNamespace(namespace); err != nil {
		return nil, err
	}

	// get
	cluster := model.NewCluster(namespace, clusterName)
	if exists, err := cluster.Select(); err != nil {
		return nil, err
	} else if exists == false {
		logger.Errorf("cluster does not exist (namespace=%s, cluster=%s)", namespace, clusterName)
		return nil, errors.New(fmt.Sprintf("Cluster not found (namespace=%s, cluster=%s)", namespace, clusterName))
	}

	return cluster, nil
}

func CreateCluster(namespace string, req *model.ClusterReq) (*model.Cluster, error) {

	// validate namespace
	if err := CheckNamespace(namespace); err != nil {
		return nil, err
	}

	clusterName := req.Name
	mcisName := clusterName

	// validate exists & clean-up cluster
	cluster := model.NewCluster(namespace, clusterName)
	if exists, err := cluster.Select(); err != nil {
		return nil, err
	} else if exists == true {
		// clean-up if "exists" & "failed-status"
		if cluster.Status.Phase == model.ClusterPhaseFailed {
			logger.Infof("clean-up cluster (namespace=%s, cluster=%s, phase=%s, reason=%s, cause=cluster is already exists) ", namespace, clusterName, cluster.Status.Phase, cluster.Status.Reason)
			_, err = DeleteCluster(namespace, clusterName)
			if err != nil {
				return nil, err
			}
		} else {
			logger.Errorf("cluster already exists (namespace=%s, cluster=%s)", namespace, clusterName)
			return nil, errors.New(fmt.Sprintf("Cluster already exists. (namespace=%s, cluster=%s)", namespace, clusterName))
		}
	}

	// start creating a cluster
	cluster = model.NewCluster(namespace, clusterName)
	cluster.NetworkCni = req.Config.Kubernetes.NetworkCni
	cluster.Label = req.Label
	cluster.InstallMonAgent = req.InstallMonAgent
	cluster.Description = req.Description

	//update phase(provisioning)
	if err := cluster.UpdatePhase(model.ClusterPhaseProvisioning); err != nil {
		return nil, err
	}

	// MCIS 존재여부 확인
	mcis := tumblebug.NewMCIS(namespace, mcisName)
	if exists, err := mcis.GET(); err != nil {
		cluster.FailReason(model.GetMCISFailedReason, err.Error())
		return nil, err
	} else if exists {
		cluster.FailReason(model.AlreadyExistMCISFailedReason, fmt.Sprintf("MCIS already exists. (namespace=%s, mcis=%s)", namespace, mcisName))
		logger.Errorf("MCIS already exists (namespace=%s, mcis=%s)", namespace, mcisName)
		return nil, errors.New(fmt.Sprintf("MCIS already exists. (namespace=%s, mcis=%s)", namespace, mcisName))
	}

	var nodeConfigInfos []NodeConfigInfo
	// control plane
	cp, err := SetNodeConfigInfos(req.ControlPlane, config.CONTROL_PLANE)
	if err != nil {
		cluster.FailReason(model.GetControlPlaneConnectionInfoFailedReason, err.Error())
		return nil, err
	}
	nodeConfigInfos = append(nodeConfigInfos, cp...)

	// worker
	wk, err := SetNodeConfigInfos(req.Worker, config.WORKER)
	if err != nil {
		cluster.FailReason(model.GetWorkerConnectionInfoFailedReason, err.Error())
		return nil, err
	}
	nodeConfigInfos = append(nodeConfigInfos, wk...)

	cIdx := 0
	wIdx := 0
	var nodes []model.Node
	var vmInfos []model.VMInfo

	for _, nodeConfigInfo := range nodeConfigInfos {
		// MCIR - 존재하면 재활용 없다면 생성 기준
		// 1. create vpc
		vpc, err := nodeConfigInfo.CreateVPC(namespace)
		if err != nil {
			cluster.FailReason(model.CreateVpcFailedReason, fmt.Sprintf("Failed to create VPC. (cause=%s)", err))
			return nil, err
		}

		// 2. create firewall
		fw, err := nodeConfigInfo.CreateFirewall(namespace)
		if err != nil {
			cluster.FailReason(model.CreateSecurityGroupFailedReason, fmt.Sprintf("Failed to create firewall. (cause=%v)", err))
			return nil, err
		}

		// 3. create sshKey
		sshKey, err := nodeConfigInfo.CreateSshKey(namespace)
		if err != nil {
			cluster.FailReason(model.CreateSSHKeyFailedReason, fmt.Sprintf("Failed to create sshkey. (cause=%v)", err))
			return nil, err
		}

		// 4. create image
		image, err := nodeConfigInfo.CreateImage(namespace)
		if err != nil {
			cluster.FailReason(model.CreateVmImageFailedReason, fmt.Sprintf("Failed to create vm-image. (cause=%v)", err))
			return nil, err
		}

		// 5. create spec
		spec, err := nodeConfigInfo.CreateSpec(namespace)
		if err != nil {
			cluster.FailReason(model.CreateVmSpecFailedReason, fmt.Sprintf("Failed to create vm-spec. (cause=%v)", err))
			return nil, err
		}

		// 6. vm
		for i := 0; i < nodeConfigInfo.Count; i++ {
			if nodeConfigInfo.Role == config.CONTROL_PLANE {
				cIdx++
			} else {
				wIdx++
			}

			vm := model.VM{
				Config:       nodeConfigInfo.Connection,
				VPC:          vpc.Name,
				Subnet:       vpc.Subnets[0].Name,
				Firewall:     []string{fw.Name},
				SSHKey:       sshKey.Name,
				Image:        image.Name,
				Spec:         spec.Name,
				UserAccount:  model.VM_USER_ACCOUNT,
				UserPassword: "",
				Description:  "",
			}

			vmInfo := model.VMInfo{
				Credential: sshKey.PrivateKey,
				Role:       nodeConfigInfo.Role,
				Csp:        nodeConfigInfo.Csp,
			}

			if nodeConfigInfo.Role == config.CONTROL_PLANE {
				vm.Name = lang.GetNodeName(config.CONTROL_PLANE, cIdx)
				if cIdx == 1 {
					vmInfo.IsCPLeader = true
					cluster.CpLeader = vm.Name
				}
			} else {
				vm.Name = lang.GetNodeName(config.WORKER, wIdx)
			}
			vmInfo.Name = vm.Name

			mcis.VMs = append(mcis.VMs, vm)
			vmInfos = append(vmInfos, vmInfo)
		}
	}

	// MCIS 생성
	mcis.Label = config.MCIS_LABEL
	mcis.InstallMonAgent = cluster.InstallMonAgent
	mcis.SystemLabel = config.MCIS_SYSTEMLABEL
	logger.Infof("start create MCIS (namespace=%s, cluster=%s, mcis=%s)", namespace, clusterName, mcisName)
	if err = mcis.POST(); err != nil {
		cluster.FailReason(model.CreateMCISFailedReason, fmt.Sprintf("Failed to create MCIS. (cause=%v)", err))
		return nil, err
	}
	logger.Infof("create MCIS OK.. (namespace=%s, cluster=%s, mcis=%s)", namespace, clusterName, mcisName)

	cpMcis := tumblebug.MCIS{}
	// 결과값 저장
	cluster.MCIS = mcisName
	for _, vm := range mcis.VMs {
		for _, vmInfo := range vmInfos {
			if vm.Name == vmInfo.Name {
				vm.Credential = vmInfo.Credential
				vm.Role = vmInfo.Role
				vm.Csp = vmInfo.Csp
				vm.IsCPLeader = vmInfo.IsCPLeader

				cpMcis.VMs = append(cpMcis.VMs, vm)
				break
			}
		}

		node := model.NewNodeVM(namespace, cluster.Name, vm)

		// insert node in store
		nodes = append(nodes, *node)
		if err := node.Insert(); err != nil {
			cluster.FailReason(model.AddNodeEntityFailedReason, fmt.Sprintf("Failed to add node entity. (cause=%v)", err))
			return nil, err
		}
	}

	// bootstrap
	logger.Infof("start k8s bootstrap (namespace=%s, cluster=%s)", namespace, clusterName)

	time.Sleep(2 * time.Second)

	eg, _ := errgroup.WithContext(context.Background())

	for _, vm := range cpMcis.VMs {
		vm := vm
		eg.Go(func() error {
			if vm.Status != config.Running || vm.PublicIP == "" {
				logger.Errorf("cannot do ssh, VM IP is not Running (namespace=%s, cluster=%s, mcis=%s, vm=%s, ip=%s, message=%s)", namespace, clusterName, mcisName, vm.Name, vm.PublicIP, vm.SystemMessage)
				return errors.New(fmt.Sprintf("Cannot do ssh, VM IP is not Running (vm=%s, ip=%s, systemMessage=%s)", vm.Name, vm.PublicIP, vm.SystemMessage))
			}
			if err := vm.ConnectionTest(); err != nil {
				return err
			}
			if err = vm.CopyScripts(cluster.NetworkCni); err != nil {
				return err
			}
			if err = vm.SetSystemd(cluster.NetworkCni); err != nil {
				return err
			}
			if err = vm.Bootstrap(); err != nil {
				return err
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		cluster.FailReason(model.SetupBoostrapFailedReason, fmt.Sprintf("Bootstrap process failed. (cause=%v)", err))
		cleanMCIS(mcis, &nodes)
		return nil, err
	}

	// init & join
	var joinCmd []string
	IPs := GetControlPlaneIPs(cpMcis.VMs)

	logger.Infof("start k8s init (namespace=%s, cluster=%s)", namespace, clusterName)
	for _, vm := range cpMcis.VMs {
		if vm.Role == config.CONTROL_PLANE && vm.IsCPLeader {

			if err := vm.InstallHAProxy(IPs); err != nil {
				cluster.FailReason(model.SetupHaproxyFailedReason, fmt.Sprintf("Failed to set up haproxy. (cause=%v)", err))
				cleanMCIS(mcis, &nodes)
				return nil, err
			}

			var clusterConfig string
			joinCmd, clusterConfig, err = vm.ControlPlaneInit(req.Config.Kubernetes)
			if err != nil {
				cluster.FailReason(model.InitControlPlaneFailedReason, fmt.Sprintf("Control-plane init. process failed. (cause=%v)", err))
				cleanMCIS(mcis, &nodes)
				return nil, err
			}
			cluster.ClusterConfig = clusterConfig

			if err = vm.InstallNetworkCNI(req.Config.Kubernetes.NetworkCni); err != nil {
				cluster.FailReason(model.SetupNetworkCNIFailedReason, fmt.Sprintf("Failed to set up network-cni. (cni=%s)", req.Config.Kubernetes.NetworkCni))
				cleanMCIS(mcis, &nodes)
				return nil, err
			}
		}
	}
	logger.Infof("end k8s init (namespace=%s, cluster=%s)", namespace, clusterName)

	logger.Infof("start k8s join (namespace=%s, cluster=%s)", namespace, clusterName)
	for _, vm := range cpMcis.VMs {
		if vm.Role == config.CONTROL_PLANE && !vm.IsCPLeader {
			if err := vm.ControlPlaneJoin(&joinCmd[0]); err != nil {
				cluster.FailReason(model.JoinControlPlaneFailedReason, fmt.Sprintf("Control-plane join process failed. (vm=%s)", vm.Name))
				cleanMCIS(mcis, &nodes)
				return nil, err
			}
		}
	}

	for _, vm := range cpMcis.VMs {
		if vm.Role == config.WORKER {
			if err := vm.WorkerJoin(&joinCmd[1]); err != nil {
				cluster.FailReason(model.JoinWorkerFailedReason, fmt.Sprintf("Worker-node join process failed. (vm=%s)", vm.Name))
				cleanMCIS(mcis, &nodes)
				return nil, err
			}
		}
	}
	logger.Infof("end k8s join (namespace=%s, cluster=%s)", namespace, clusterName)

	logger.Infof("start add node labels (namespace=%s, cluster=%s)", namespace, clusterName)
	for _, vm := range cpMcis.VMs {
		if err := vm.AddNodeLabels(); err != nil {
			logger.Warnf("failed to add node labels (namespace=%s, cluster=%s, vm=%s, cause=%s)", namespace, clusterName, vm.Name, err)
		}
	}
	logger.Infof("end add node labels (namespace=%s, cluster=%s)", namespace, clusterName)

	cluster.UpdatePhase(model.ClusterPhaseProvisioned)

	nodeList, _ := updateNodesCreatedTime(namespace, clusterName, &nodes)
	cluster.Nodes = nodeList

	return cluster, nil
}

func DeleteCluster(namespace string, clusterName string) (*model.Status, error) {

	// validate namespace
	if err := CheckNamespace(namespace); err != nil {
		return nil, err
	}
	// validate exists
	cluster := model.NewCluster(namespace, clusterName)
	exists, err := cluster.Select()
	if err != nil {
		return nil, err
	} else if exists == false {
		logger.Errorf("cluster not found (namespace=%s, cluster=%s)", namespace, clusterName)
		return nil, errors.New(fmt.Sprintf("Cluster not found (namespace=%s, cluster=%s)", namespace, clusterName))
	}

	mcisName := clusterName
	status := model.NewStatus()
	status.Code = model.STATUS_UNKNOWN

	logger.Infof("start delete Cluster (namespace=%s, cluster=%s)", namespace, clusterName)

	// 0. set stauts
	cluster.UpdatePhase(model.ClusterPhaseDeleting)

	// 1.delete MCIS
	mcis := tumblebug.NewMCIS(namespace, mcisName)
	exist, err := mcis.GET()
	if err != nil {
		return nil, err
	} else if exist {
		if err = deleteMCIS(mcis); err != nil {
			return nil, err
		}
	}

	// 2.delete cluster-entity
	if err := cluster.Delete(); err != nil {
		status.Message = fmt.Sprintf("Failed to delete cluster (namespace=%s, cluster=%s)", namespace, clusterName)
		return nil, err
	}

	status.Code = model.STATUS_SUCCESS
	status.Message = fmt.Sprintf("cluster %s has been deleted", mcisName)
	return status, nil
}

func cleanMCIS(mcis *tumblebug.MCIS, nodesModel *[]model.Node) error {

	logger.Infof("clean MCIS (mcis=%s)", mcis.Name)
	for _, nd := range *nodesModel {
		if err := nd.Delete(); err != nil {
			logger.Warnf("failed to clean MCIS (op=delete, mcis=%s, cause=%V)", mcis.Name, err)
		}
		nd.Credential = ""
		nd.PublicIP = ""
		if err := nd.Insert(); err != nil {
			logger.Warnf("failed to clean MCIS (op=insert, mcis=%s, cause=%V)", mcis.Name, err)
		}
	}

	return deleteMCIS(mcis)
}

func deleteMCIS(mcis *tumblebug.MCIS) error {

	mcisName := mcis.Name

	logger.Infof("terminate MCIS (mcis=%s)", mcisName)
	if err := mcis.TERMINATE(); err != nil {
		return errors.New(fmt.Sprintf("Failed to terminate MCIS (mcis=%s)", mcis.Name))
	}
	time.Sleep(5 * time.Second)

	logger.Infof("delete MCIS (mcis=%s)", mcisName)
	if err := mcis.DELETE(); err != nil {
		if strings.Contains(err.Error(), "Deletion is not allowed") {
			logger.Infof("refine mcis (mcis=%s)", mcisName)
			if err = mcis.REFINE(); err != nil {
				return errors.New(fmt.Sprintf("Failed to refine MCIS (cause=%v)", err))
			}
			logger.Infof("delete MCIS (mcis=%s)", mcisName)
			if err = mcis.DELETE(); err != nil {
				return errors.New(fmt.Sprintf("Failed to delete MCIS (cause=%v)", err))
			}
		} else {
			return errors.New(fmt.Sprintf("Failed to delete MCIS (cause=%v)", err))
		}
	}

	logger.Infof("delete MCIS OK.. (mcis=%s)", mcisName)
	return nil
}

func updateNodesCreatedTime(namespace string, clusterName string, nodesModel *[]model.Node) ([]model.Node, error) {
	for _, node := range *nodesModel {
		node.CreatedTime = lang.GetNowUTC()
		if err := node.Insert(); err != nil {
			logger.Warnf("failed to update node's createdtime (mcis=%s, cause=%v)", clusterName, err)
		}
	}

	var nodes []model.Node
	nodeList := model.NewNodeList(namespace, clusterName)
	err := nodeList.SelectList()
	if err != nil {
		logger.Warnf("failed to select node list (mcis=%s, cause=%v)", clusterName, err)
		nodes = *nodesModel
	}
	nodes = nodeList.Items

	return nodes, nil
}
