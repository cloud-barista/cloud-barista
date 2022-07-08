package service

import (
	//"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloud-barista/cb-mcks/src/core/app"
	"github.com/cloud-barista/cb-mcks/src/core/model"
	"github.com/cloud-barista/cb-mcks/src/core/provision"
	"github.com/cloud-barista/cb-mcks/src/core/tumblebug"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"

	logger "github.com/sirupsen/logrus"
)

/* get nodes */
func ListNode(namespace string, clusterName string) (*model.NodeList, error) {
	err := verifyNamespace(namespace)
	if err != nil {
		return nil, err
	}

	nodeList := &model.NodeList{
		ListModel: model.ListModel{Kind: app.KIND_NODE_LIST},
		Items:     []*model.Node{},
	}

	cluster := model.NewCluster(namespace, clusterName)
	if exist, err := cluster.Select(); err != nil {
		return nil, err
	} else if exist == true {
		nodeList.Items = cluster.Nodes
	}
	return nodeList, nil
}

/* get a node */
func GetNode(namespace string, clusterName string, nodeName string) (*model.Node, error) {
	err := verifyNamespace(namespace)
	if err != nil {
		return nil, err
	}

	cluster := model.NewCluster(namespace, clusterName)
	if exists, err := cluster.Select(); err != nil {
		return nil, err
	} else if exists {
		for _, node := range cluster.Nodes {
			if node.Name == nodeName {
				return node, nil
			}
		}
	}

	return nil, errors.New(fmt.Sprintf("Could not be found a node '%s' (namespace=%s, cluster=%s)", nodeName, namespace, clusterName))
}

/* add a node */
func AddNode(namespace string, clusterName string, req *app.NodeReq) (*model.NodeList, error) {

	// validate namespace
	if err := verifyNamespace(namespace); err != nil {
		return nil, err
	}

	// get a cluster-entity
	cluster := model.NewCluster(namespace, clusterName)
	if exists, err := cluster.Select(); err != nil {
		return nil, err
	} else if exists == false {
		return nil, errors.New(fmt.Sprintf("Could not be found a cluster '%s'. (namespace=%s)", clusterName, namespace))
	} else if cluster.Status.Phase != model.ClusterPhaseProvisioned {
		return nil, errors.New(fmt.Sprintf("Unable to add a node. status is '%s'.", cluster.Status.Phase))
	}

	// get a MCIS
	mcis := tumblebug.NewMCIS(namespace, cluster.MCIS)
	if exists, err := mcis.GET(); err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.New(fmt.Sprintf("Can't be found a MCIS '%s'.", cluster.MCIS))
	}
	logger.Infof("[%s.%s] The inquiry has been completed..", namespace, clusterName)

	mcisName := cluster.MCIS

	// get a provisioner
	provisioner := provision.NewProvisioner(cluster)

	// get join command
	workerJoinCmd, err := provisioner.NewWorkerJoinCommand()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to get join-command (cause='%v')", err))
	}
	logger.Infof("[%s.%s] Worker join-command inquiry has been completed. (command=%s)", namespace, clusterName, workerJoinCmd)

	// create a MCIR & MCIS-vm
	idx := cluster.NextNodeIndex(app.WORKER)
	vms := []tumblebug.VM{}
	for _, worker := range req.Worker {
		mcir := NewMCIR(namespace, app.WORKER, worker)
		reason, msg := mcir.CreateIfNotExist()
		if reason != "" {
			return nil, errors.New(msg)
		} else {
			for i := 0; i < mcir.vmCount; i++ {
				name := lang.GenerateNewNodeName(string(app.WORKER), idx)
				vm := mcir.NewVM(namespace, name, mcisName)
				if err := vm.POST(); err != nil {
					cleanUpNodes(*provisioner)
					return nil, err
				}
				vms = append(vms, vm)
				provisioner.AppendWorkerNodeMachine(name, mcir.csp, mcir.region, mcir.zone, mcir.credential)
				idx = idx + 1
			}
		}
	}
	logger.Infof("[%s.%s] MCIS(vm) creation has been completed. (len=%d)", namespace, clusterName, len(vms))

	// save nodes metadata
	if nodes, err := provisioner.BindVM(vms); err != nil {
		return nil, err
	} else {
		cluster.Nodes = append(cluster.Nodes, nodes...)
		if err := cluster.PutStore(); err != nil {
			cleanUpNodes(*provisioner)
			return nil, errors.New(fmt.Sprintf("Failed to add node entity. (cause='%v')", err))
		}
	}

	// kubernetes provisioning : bootstrap
	time.Sleep(2 * time.Second)
	if err := provisioner.Bootstrap(); err != nil {
		cleanUpNodes(*provisioner)
		return nil, errors.New(fmt.Sprintf("Bootstrap failed. (cause='%v')", err))
	}
	logger.Infof("[%s.%s] Bootstrap has been completed.", namespace, clusterName)

	// kubernetes provisioning : worker node join
	for _, machine := range provisioner.WorkerNodeMachines {
		if err := machine.JoinWorker(&workerJoinCmd); err != nil {
			cleanUpNodes(*provisioner)
			return nil, errors.New(fmt.Sprintf("Fail to worker-node join. (node=%s)", machine.Name))
		}
	}
	logger.Infof("[%s.%s] Woker-nodes join has been completed.", namespace, clusterName)

	// assign node labels (topology.cloud-barista.github.io/csp , topology.kubernetes.io/region, topology.kubernetes.io/zone)
	if err = provisioner.AssignNodeLabelAnnotation(); err != nil {
		logger.Warnf("[%s.%s] Failed to assign node labels (cause='%v')", namespace, clusterName, err)
	} else {
		logger.Infof("[%s.%s] Node label assignment has been completed.", namespace, clusterName)
	}

	// save nodes metadata & update status
	for _, node := range cluster.Nodes {
		node.CreatedTime = lang.GetNowUTC()
	}
	if err := cluster.PutStore(); err != nil {
		cleanUpNodes(*provisioner)
		return nil, errors.New(fmt.Sprintf("Failed to add node entity. (cause='%v')", err))
	}
	logger.Infof("[%s.%s] Nodes creation has been completed.", namespace, clusterName)

	nodes := model.NewNodeList(namespace, clusterName)
	nodes.Items = cluster.Nodes
	return nodes, nil
}

/* remove a node */
func RemoveNode(namespace string, clusterName string, nodeName string) (*app.Status, error) {

	//validate
	if err := verifyNamespace(namespace); err != nil {
		return nil, err
	}

	// get a cluster-entity
	cluster := model.NewCluster(namespace, clusterName)
	if exists, err := cluster.Select(); err != nil {
		return nil, err
	} else if exists == false {
		return nil, errors.New(fmt.Sprintf("Could not be found a cluster. (namespace=%s, cluster=%s)", namespace, clusterName))
	} else if cluster.Status.Phase != model.ClusterPhaseProvisioned {
		return nil, errors.New(fmt.Sprintf("Unable to remove a node. status is '%s'.", cluster.Status.Phase))
	}

	// validate exists
	if nodeName == cluster.CpLeader {
		return nil, errors.New("Could not be delete a control-plane leader node.")
	}
	if exists := cluster.ExistsNode(nodeName); !exists {
		return app.NewStatus(app.STATUS_NOTFOUND, fmt.Sprintf("Could not be found a node-entity '%s'", nodeName)), nil
	}
	logger.Infof("[%s.%s] The inquiry has been completed..", namespace, clusterName)

	// get a provisioner
	provisioner := provision.NewProvisioner(cluster)
	// delete node (kubernetes) & vm (mcis)
	if err := provisioner.DrainAndDeleteNode(nodeName); err != nil {
		return nil, err
	}
	// delete a node-entity
	if err := cluster.DeleteNode(nodeName); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to delete a cluster-entity. (cause='%v')", err))
	}

	logger.Infof("[%s.%s] Node deletinn has been completed. (node=%s)", namespace, clusterName, nodeName)
	return app.NewStatus(app.STATUS_SUCCESS, fmt.Sprintf("Node '%s' has been deleted", nodeName)), nil
}

/* clean-up nodes (with VMs) & update a node-entities */
func cleanUpNodes(provisioner provision.Provisioner) {

	for _, machine := range provisioner.GetMachinesAll() {
		nodeName := machine.Name
		existNode := false
		for _, node := range provisioner.Cluster.Nodes {
			if node.Name == nodeName {
				node.Credential = ""
				node.PublicIP = ""
				existNode = true
				break
			}
		}
		if existNode {
			if err := provisioner.DrainAndDeleteNode(nodeName); err != nil {
				logger.Warnf("[%s.%s] %s", provisioner.Cluster.Namespace, provisioner.Cluster.Name, err.Error())
			}
		}
	}
	if err := provisioner.Cluster.PutStore(); err != nil {
		logger.Warnf("[%s.%s] Failed to update a cluster-entity. (cause='%v')", provisioner.Cluster.Namespace, provisioner.Cluster.Name, err)
	}
	logger.Infof("[%s.%s] Garbage data has been cleaned.", provisioner.Cluster.Namespace, provisioner.Cluster.Name)
}
