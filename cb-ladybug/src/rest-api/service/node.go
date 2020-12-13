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
)

func ListNode(namespace string, clusterName string) (*model.NodeList, error) {
	nodes := model.NewNodeList()
	result, err := nodes.SelectList(namespace, clusterName, nodes)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetNode(namespace string, clusterName string, nodeName string) (*model.Node, error) {
	node := model.NewNodeDef(nodeName)
	result, err := node.Select(namespace, clusterName, node)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func AddNode(namespace string, clusterName string, req *model.NodeReq) (*model.NodeList, error) {

	//TODO [update/hard-coding] connection config
	csp := config.CSP_GCP
	if strings.Contains(namespace, "aws") {
		csp = config.CSP_AWS
	}
	//host user account
	account := GetUserAccount(csp)

	// get join command
	cpNode, err := getCPNode(namespace, clusterName)
	if err != nil {
		return nil, errors.New("control-plane node not found")
	}
	workerJoinCmd, err := getWorkerJoinCmdForAddNode(account, cpNode)
	if err != nil {
		return nil, errors.New("get join command error")
	}

	mcisName := clusterName
	mcis := tumblebug.NewMCIS(namespace, mcisName)

	exists, err := mcis.GET()
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("MCIS not found")
	}

	vpcName := fmt.Sprintf("%s-vpc", clusterName)
	firewallName := fmt.Sprintf("%s-allow-external", clusterName)
	sshkeyName := fmt.Sprintf("%s-sshkey", clusterName)
	imageName := fmt.Sprintf("%s-Ubuntu1804", req.Config)
	specName := fmt.Sprintf("%s-spec", clusterName)

	// vpc
	fmt.Println(fmt.Sprintf("start create vpc (name=%s)", vpcName))
	vpc := tumblebug.NewVPC(namespace, vpcName, req.Config)
	exists, e := vpc.GET()
	if e != nil {
		return nil, e
	}
	if exists {
		fmt.Println(fmt.Sprintf("reuse vpc (name=%s, cause='already exists')", vpcName))
	} else {
		if e = vpc.POST(); e != nil {
			return nil, e
		}
		fmt.Println(fmt.Sprintf("create vpc OK.. (name=%s)", vpcName))
	}

	// firewall
	fmt.Println(fmt.Sprintf("start create firewall (name=%s)", firewallName))
	fw := tumblebug.NewFirewall(namespace, firewallName, req.Config)
	fw.VPCId = vpcName
	exists, e = fw.GET()
	if e != nil {
		return nil, e
	}
	if exists {
		fmt.Println(fmt.Sprintf("reuse firewall (name=%s, cause='already exists')", firewallName))
	} else {
		if e = fw.POST(); e != nil {
			return nil, e
		}
		fmt.Println(fmt.Sprintf("create firewall OK.. (name=%s)", firewallName))
	}

	// sshKey
	fmt.Println(fmt.Sprintf("start create ssh key (name=%s)", sshkeyName))
	sshKey := tumblebug.NewSSHKey(namespace, sshkeyName, req.Config)
	sshKey.Username = "cb-cluster"
	exists, e = sshKey.GET()
	if e != nil {
		return nil, e
	}
	if exists {
		fmt.Println(fmt.Sprintf("reuse ssh key (name=%s, cause='already exists')", sshkeyName))
	} else {
		if e = sshKey.POST(); e != nil {
			return nil, e
		}
		fmt.Println(fmt.Sprintf("create ssh key OK.. (name=%s)", sshkeyName))
	}

	// image
	fmt.Println(fmt.Sprintf("start create image (name=%s)", imageName))
	// get image id
	imageId, e := GetVmImageId(csp, req.Config)
	if e != nil {
		return nil, e
	}

	image := tumblebug.NewImage(namespace, imageName, req.Config)
	image.CspImageId = imageId
	exists, e = image.GET()
	if e != nil {
		return nil, e
	}
	if exists {
		fmt.Println(fmt.Sprintf("reuse image (name=%s, cause='already exists')", imageName))
	} else {
		if e = image.POST(); e != nil {
			return nil, e
		}
		fmt.Println(fmt.Sprintf("create image OK.. (name=%s)", imageName))
	}

	// spec
	fmt.Println(fmt.Sprintf("start create worker node spec (name=%s)", specName))
	spec := tumblebug.NewSpec(namespace, specName, req.Config)
	spec.CspSpecName = req.WorkerNodeSpec
	spec.Role = config.WORKER
	exists, e = spec.GET()
	if e != nil {
		return nil, e
	}
	if exists {
		fmt.Println(fmt.Sprintf("reuse worker node spec (name=%s, cause='already exists')", specName))
	} else {
		if e = spec.POST(); e != nil {
			return nil, e
		}
		fmt.Println(fmt.Sprintf("create worker node spec OK.. (name=%s)", specName))
	}

	// vm
	var VMs []model.VM
	for i := 0; i < req.WorkerNodeCount; i++ {
		vm := tumblebug.NewTVm(namespace, mcisName)
		vm.VM = model.VM{
			Name:         lang.GetNodeName(clusterName, spec.Role),
			Config:       req.Config,
			VPC:          vpc.Name,
			Subnet:       vpc.Subnets[0].Name,
			Firewall:     []string{fw.Name},
			SSHKey:       sshKey.Name,
			Image:        image.Name,
			Spec:         spec.Name,
			UserAccount:  account,
			UserPassword: "",
			Description:  "",
			Credential:   sshKey.PrivateKey,
			Role:         spec.Role,
		}

		// vm 생성
		fmt.Println(fmt.Sprintf("start create VM (mcisname=%s, nodename=%s)", mcisName, vm.VM.Name))
		err := vm.POST()
		if err != nil {
			fmt.Println(fmt.Sprintf("create VM error (mcisname=%s, nodename=%s)", mcisName, vm.VM.Name))
			return nil, err
		}
		VMs = append(VMs, vm.VM)
		fmt.Println(fmt.Sprintf("create VM OK.. (mcisname=%s, nodename=%s)", mcisName, vm.VM.Name))
	}

	var wg sync.WaitGroup
	c := make(chan error)
	wg.Add(len(VMs))

	fmt.Println("start connect VMs")
	for _, vm := range VMs {
		go func(vm model.VM) {
			defer wg.Done()
			sshInfo := ssh.SSHInfo{
				UserName:   account,
				PrivateKey: []byte(vm.Credential),
				ServerPort: fmt.Sprintf("%s:22", vm.PublicIP),
			}

			_ = vm.ConnectionTest(&sshInfo, &vm)
			err := vm.CopyScripts(&sshInfo, &vm)
			if err != nil {
				c <- err
			}
			bootstrapResult, err := vm.Bootstrap(&sshInfo)
			if err != nil {
				c <- err
			}
			if !bootstrapResult {
				c <- errors.New(vm.Name + " bootstrap failed")
			}
			result, err := vm.WorkerJoinForAddNode(&sshInfo, &workerJoinCmd)
			if err != nil {
				c <- err
			}
			if !result {
				c <- errors.New(vm.Name + " join failed")
			}
		}(vm)
	}

	go func() {
		wg.Wait()
		close(c)
		fmt.Println("end connect VMs")
	}()

	for err := range c {
		if err != nil {
			common.CBLog.Error(" worker join error == ", err)
			return nil, err
		}
	}

	// insert store
	nodes := model.NewNodeList()
	for _, vm := range VMs {
		node := model.NewNode(vm)
		result, err := node.Insert(namespace, clusterName, node)
		if err != nil {
			return nil, err
		}
		nodes.Items = append(nodes.Items, *result)
	}

	return nodes, nil
}

func RemoveNode(namespace string, clusterName string, nodeName string) (*model.Status, error) {
	status := model.NewStatus()
	status.Code = model.STATUS_UNKNOWN

	cpNode, err := getCPNode(namespace, clusterName)
	if err != nil {
		status.Message = "control-plane node not found"
		return status, err
	}

	var userAccount string
	var hostName string
	if strings.Contains(namespace, "gcp") {
		userAccount = "cb-user"
		hostName = nodeName
	} else {
		userAccount = "ubuntu"

		hostName, err = getAWSHostName(namespace, clusterName, nodeName, userAccount)
		if err != nil {
			status.Message = "get aws node name error"
			return status, err
		}
	}

	// drain node
	sshInfo := ssh.SSHInfo{
		UserName:   userAccount,
		PrivateKey: []byte(cpNode.Credential),
		ServerPort: fmt.Sprintf("%s:22", cpNode.PublicIP),
	}
	cmd := fmt.Sprintf("sudo kubectl drain %s --kubeconfig=/etc/kubernetes/admin.conf --ignore-daemonsets", hostName)
	result, err := ssh.SSHRun(sshInfo, cmd)
	if err != nil {
		status.Message = "kubectl drain failed"
		return status, err
	}
	if strings.Contains(result, fmt.Sprintf("node/%s drained", hostName)) || strings.Contains(result, fmt.Sprintf("node/%s evicted", hostName)) {
		fmt.Println("drain node success")
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
		fmt.Println("delete node success")
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
	node := model.NewNodeDef(nodeName)
	if err := node.Delete(namespace, clusterName, node); err != nil {
		status.Message = err.Error()
		return status, nil
	}

	status.Code = model.STATUS_SUCCESS
	status.Message = "success"

	return status, nil
}

func getCPNode(namespace string, clusterName string) (*model.Node, error) {
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
		if node.Role == config.CONTROL_PLANE {
			cpNode = node
			break
		}
	}

	return cpNode, nil
}

func getAWSHostName(namespace string, clusterName string, nodeName string, userAccount string) (string, error) {
	key := lang.GetStoreNodeKey(namespace, clusterName, "")
	keyValues, err := common.CBStore.GetList(key, true)
	if err != nil {
		return "", err
	}
	if keyValues == nil {
		return "", errors.New(fmt.Sprintf("%s not found", key))
	}
	wNode := &model.Node{}
	for _, keyValue := range keyValues {
		node := &model.Node{}
		json.Unmarshal([]byte(keyValue.Value), &node)
		if node.Name == nodeName {
			wNode = node
			break
		}
	}

	sshInfo := ssh.SSHInfo{
		UserName:   userAccount,
		PrivateKey: []byte(wNode.Credential),
		ServerPort: fmt.Sprintf("%s:22", wNode.PublicIP),
	}
	cmd := "/bin/hostname"
	result, err := ssh.SSHRun(sshInfo, cmd)
	if err != nil {
		return "", err
	}

	return result, nil
}

func getWorkerJoinCmdForAddNode(userAccount string, cpNode *model.Node) (string, error) {
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
	// privateIp -> publicIP
	// if strings.Contains(result, "kubeadm join") {
	// 	ip := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	// 	ipRegEx, _ := regexp.Compile(fmt.Sprintf("%s\\.%s\\.%s\\.%s", ip, ip, ip, ip))
	// 	if ipRegEx.MatchString(result) {
	// 		res := ipRegEx.FindString(result)
	// 		result = strings.Replace(result, res, cpNode.PublicIP, 1)
	// 	}
	// } else {
	// 	return "", errors.New("get join command error")
	// }

	return result, nil
}
