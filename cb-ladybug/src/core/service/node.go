package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/cloud-barista/cb-ladybug/src/core/common"
	"github.com/cloud-barista/cb-ladybug/src/core/model"
	"github.com/cloud-barista/cb-ladybug/src/core/model/tumblebug"
	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	"github.com/cloud-barista/cb-ladybug/src/utils/lang"

	ssh "github.com/cloud-barista/cb-spider/cloud-control-manager/vm-ssh"
	logger "github.com/sirupsen/logrus"
)

func ListNode(namespace string, clusterName string) (*model.NodeList, error) {
	nodes := model.NewNodeList(namespace, clusterName)
	err := nodes.SelectList()
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func GetNode(namespace string, clusterName string, nodeName string) (*model.Node, error) {
	node := model.NewNode(namespace, clusterName, nodeName)
	err := node.Select()
	if err != nil {
		return nil, err
	}

	return node, nil
}

func AddNode(namespace string, clusterName string, req *model.NodeReq) (*model.NodeList, error) {

	mcisName := clusterName
	mcis := tumblebug.NewMCIS(namespace, mcisName)

	exists, err := mcis.GET()
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("MCIS not found")
	}

	// get join command
	workerJoinCmd, err := getWorkerJoinCmdForAddNode(namespace, clusterName)
	if err != nil {
		return nil, errors.New("join command cannot get")
	}
	networkCni := getClusterNetworkCNI(namespace, clusterName)

	var nodeConfigInfos []NodeConfigInfo
	// worker
	wk, err := SetNodeConfigInfos(req.Worker, config.WORKER)
	if err != nil {
		return nil, err
	}
	nodeConfigInfos = append(nodeConfigInfos, wk...)

	cIdx := 0
	wIdx := 0
	maxCIdx, maxWIdx := getMaxIdx(namespace, clusterName)
	var TVMs []tumblebug.TVM

	for _, nodeConfigInfo := range nodeConfigInfos {
		// MCIR - 존재하면 재활용 없다면 생성 기준
		// 1. create vpc
		vpc, err := nodeConfigInfo.CreateVPC(namespace)
		if err != nil {
			return nil, err
		}

		// 2. create firewall
		fw, err := nodeConfigInfo.CreateFirewall(namespace)
		if err != nil {
			return nil, err
		}

		// 3. create sshKey
		sshKey, err := nodeConfigInfo.CreateSshKey(namespace)
		if err != nil {
			return nil, err
		}

		// 4. create image
		image, err := nodeConfigInfo.CreateImage(namespace)
		if err != nil {
			return nil, err
		}

		// 5. create spec
		spec, err := nodeConfigInfo.CreateSpec(namespace)
		if err != nil {
			return nil, err
		}

		// 6. vm
		for i := 0; i < nodeConfigInfo.Count; i++ {
			if nodeConfigInfo.Role == config.CONTROL_PLANE {
				cIdx++
			} else {
				wIdx++
			}
			tvm := tumblebug.NewTVm(namespace, mcisName)
			tvm.VM = model.VM{
				Config:       nodeConfigInfo.Connection,
				VPC:          vpc.Name,
				Subnet:       vpc.Subnets[0].Name,
				Firewall:     []string{fw.Name},
				SSHKey:       sshKey.Name,
				Image:        image.Name,
				Spec:         spec.Name,
				UserAccount:  nodeConfigInfo.Account,
				UserPassword: "",
				Description:  "",
				Credential:   sshKey.PrivateKey,
				Role:         nodeConfigInfo.Role,
				Csp:          nodeConfigInfo.Csp,
			}

			if nodeConfigInfo.Role == config.CONTROL_PLANE {
				tvm.VM.Name = lang.GetNodeName(clusterName, config.CONTROL_PLANE, maxCIdx+cIdx)
			} else {
				tvm.VM.Name = lang.GetNodeName(clusterName, config.WORKER, maxWIdx+wIdx)
			}

			// vm 생성
			logger.Infof("start create VM (mcisname=%s, nodename=%s)", mcisName, tvm.VM.Name)
			err := tvm.POST()
			if err != nil {
				logger.Warnf("create VM error (mcisname=%s, nodename=%s)", mcisName, tvm.VM.Name)
				return nil, err
			}
			logger.Infof("create VM OK.. (mcisname=%s, nodename=%s)", mcisName, tvm.VM.Name)

			TVMs = append(TVMs, *tvm)
		}
	}

	var wg sync.WaitGroup
	c := make(chan error)
	wg.Add(len(TVMs))

	logger.Infoln("start connect VMs")
	for _, tvm := range TVMs {
		go func(vm model.VM) {
			defer wg.Done()
			sshInfo := ssh.SSHInfo{
				UserName:   GetUserAccount(vm.Csp),
				PrivateKey: []byte(vm.Credential),
				ServerPort: fmt.Sprintf("%s:22", vm.PublicIP),
			}

			_ = vm.ConnectionTest(&sshInfo)
			err := vm.CopyScripts(&sshInfo, networkCni)
			if err != nil {
				c <- err
			}

			logger.Infof("set systemd service (vm=%s)", vm.Name)
			err = vm.SetSystemd(&sshInfo, networkCni)
			if err != nil {
				c <- err
			}

			logger.Infof("bootstrap (vm=%s)", vm.Name)
			err = vm.Bootstrap(&sshInfo)
			if err != nil {
				c <- err
			}

			logger.Infof("join (vm=%s)", vm.Name)
			err = vm.WorkerJoin(&sshInfo, &workerJoinCmd)
			if err != nil {
				c <- err
			}
		}(tvm.VM)
	}

	go func() {
		wg.Wait()
		close(c)
		logger.Infoln("end connect VMs")
	}()

	for err := range c {
		if err != nil {
			logger.Warnf("worker join error (cause=%v)", err)
			return nil, err
		}
	}

	// insert store
	nodes := model.NewNodeList(namespace, clusterName)
	for _, vm := range TVMs {
		node := model.NewNodeVM(namespace, clusterName, vm.VM)
		node.UId = lang.GetUid()
		err := node.Insert()
		if err != nil {
			return nil, err
		}
		nodes.Items = append(nodes.Items, *node)
	}

	return nodes, nil
}

func RemoveNode(namespace string, clusterName string, nodeName string) (*model.Status, error) {
	status := model.NewStatus()
	status.Code = model.STATUS_UNKNOWN

	cpNode, err := getCPLeaderNode(namespace, clusterName)
	if err != nil {
		status.Message = "control-plane node not found"
		return status, err
	}

	hostName, err := getHostName(namespace, clusterName, nodeName)
	if err != nil {
		status.Message = "get node name error"
		return status, err
	}

	// drain node
	userAccount := GetUserAccount(cpNode.Csp)
	sshInfo := ssh.SSHInfo{
		UserName:   userAccount,
		PrivateKey: []byte(cpNode.Credential),
		ServerPort: fmt.Sprintf("%s:22", cpNode.PublicIP),
	}
	cmd := fmt.Sprintf("sudo kubectl drain %s --kubeconfig=/etc/kubernetes/admin.conf --ignore-daemonsets --force --delete-local-data", hostName)
	result, err := ssh.SSHRun(sshInfo, cmd)
	if err != nil {
		status.Message = "kubectl drain failed"
		return status, err
	}
	if strings.Contains(result, fmt.Sprintf("node/%s drained", hostName)) || strings.Contains(result, fmt.Sprintf("node/%s evicted", hostName)) {
		logger.Infoln("drain node success")
	} else {
		status.Message = "kubectl drain failed"
		return status, err
	}

	// delete node
	cmd = fmt.Sprintf("sudo kubectl delete node %s --kubeconfig=/etc/kubernetes/admin.conf", hostName)
	result, err = ssh.SSHRun(sshInfo, cmd)
	if err != nil {
		status.Message = "kubectl delete node failed"
		return status, err
	}
	if strings.Contains(result, "deleted") {
		logger.Infoln("delete node success")
	} else {
		status.Message = "kubectl delete node failed"
		return status, errors.New("kubectl delete node failed")
	}

	// delete vm
	vm := tumblebug.NewTVm(namespace, clusterName)
	vm.VM.Name = nodeName
	err = vm.DELETE()
	if err != nil {
		status.Message = "delete vm failed"
		return status, err
	}

	// delete node in store
	node := model.NewNode(namespace, clusterName, nodeName)
	if err := node.Delete(); err != nil {
		status.Message = err.Error()
		return status, nil
	}

	status.Code = model.STATUS_SUCCESS
	status.Message = "success"

	return status, nil
}

func getCluster(namespace string, clusterName string) (*model.Cluster, error) {
	key := lang.GetStoreClusterKey(namespace, clusterName)
	keyValue, err := common.CBStore.Get(key)
	if err != nil {
		return nil, err
	}
	if keyValue == nil {
		return nil, errors.New(fmt.Sprintf("%s not found", key))
	}
	cluster := &model.Cluster{}
	json.Unmarshal([]byte(keyValue.Value), cluster)

	return cluster, nil
}

func getCPLeaderNode(namespace string, clusterName string) (*model.Node, error) {
	cluster, err := getCluster(namespace, clusterName)
	if err != nil {
		return nil, errors.New("cluster info not found")
	}
	CpLeader := cluster.CpLeader

	key := lang.GetStoreNodeKey(namespace, clusterName, "")
	keyValues, err := common.CBStore.GetList(key, true)
	if err != nil {
		return nil, err
	}
	if keyValues == nil {
		return nil, errors.New(fmt.Sprintf("%s not found", key))
	}
	cpNode := &model.Node{}
	for _, keyValue := range keyValues {
		node := &model.Node{}
		json.Unmarshal([]byte(keyValue.Value), &node)
		if node.Name == CpLeader {
			cpNode = node
			break
		}
	}

	return cpNode, nil
}

func getClusterNetworkCNI(namespace string, clusterName string) string {
	cluster, err := getCluster(namespace, clusterName)
	if err != nil {
		return ""
	}

	return cluster.NetworkCni
}

func getHostName(namespace string, clusterName string, nodeName string) (string, error) {
	key := lang.GetStoreNodeKey(namespace, clusterName, "")
	keyValues, err := common.CBStore.GetList(key, true)
	if err != nil {
		return "", err
	}
	if keyValues == nil {
		return "", errors.New(fmt.Sprintf("%s not found", key))
	}
	dNode := &model.Node{}
	for _, keyValue := range keyValues {
		node := &model.Node{}
		json.Unmarshal([]byte(keyValue.Value), &node)
		if node.Name == nodeName {
			dNode = node
			break
		}
	}

	userAccount := GetUserAccount(dNode.Csp)
	sshInfo := ssh.SSHInfo{
		UserName:   userAccount,
		PrivateKey: []byte(dNode.Credential),
		ServerPort: fmt.Sprintf("%s:22", dNode.PublicIP),
	}
	cmd := "/bin/hostname"
	hostName, err := ssh.SSHRun(sshInfo, cmd)
	if err != nil {
		return "", err
	}
	return hostName, nil
}

func getWorkerJoinCmdForAddNode(namespace string, clusterName string) (string, error) {
	cpNode, err := getCPLeaderNode(namespace, clusterName)
	if err != nil {
		return "", errors.New("control-plane node not found")
	}
	userAccount := GetUserAccount(cpNode.Csp)
	sshInfo := ssh.SSHInfo{
		UserName:   userAccount,
		PrivateKey: []byte(cpNode.Credential),
		ServerPort: fmt.Sprintf("%s:22", cpNode.PublicIP),
	}
	cmd := "sudo kubeadm token create --print-join-command"
	result, err := ssh.SSHRun(sshInfo, cmd)
	if err != nil {
		return "", err
	}
	return result, nil
}

func getMaxIdx(namespace string, clusterName string) (maxCpIdx int, maxWkIdx int) {
	maxCpIdx = 0
	maxWkIdx = 0

	nodes := model.NewNodeList(namespace, clusterName)
	err := nodes.SelectList()
	if err != nil {
		return
	}

	var arrCp, arrWk []int
	for _, node := range nodes.Items {
		slice := strings.Split(node.Name, "-")
		role := len(slice) - 3
		idx := len(slice) - 2

		if slice[role] == "c" {
			arrCp = append(arrCp, lang.GetIdxToInt(slice[idx]))
		} else if slice[role] == "w" {
			arrWk = append(arrWk, lang.GetIdxToInt(slice[idx]))
		}
	}
	fmt.Println(maxCpIdx, maxWkIdx)
	maxCpIdx = lang.GetMaxNumber(arrCp)
	maxWkIdx = lang.GetMaxNumber(arrWk)
	return
}
