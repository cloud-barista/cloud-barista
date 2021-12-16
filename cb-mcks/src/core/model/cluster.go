package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cloud-barista/cb-mcks/src/core/common"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"

	logger "github.com/sirupsen/logrus"
)

type ClusterPhase string
type ClusterReason string

const (
	ClusterPhasePending      = ClusterPhase("Pending")
	ClusterPhaseProvisioning = ClusterPhase("Provisioning")
	ClusterPhaseProvisioned  = ClusterPhase("Provisioned")
	ClusterPhaseFailed       = ClusterPhase("Failed")
	ClusterPhaseDeleting     = ClusterPhase("Deleting")

	GetMCISFailedReason                       = ClusterReason("GetMCISFailedReason")
	AlreadyExistMCISFailedReason              = ClusterReason("AlreadyExistMCISFailedReason")
	CreateMCISFailedReason                    = ClusterReason("CreateMCISFailedReason")
	GetControlPlaneConnectionInfoFailedReason = ClusterReason("GetControlPlaneConnectionInfoFailedReason")
	GetWorkerConnectionInfoFailedReason       = ClusterReason("GetWorkerConnectionInfoFailedReason")
	CreateVpcFailedReason                     = ClusterReason("CreateVpcFailedReason")
	CreateSecurityGroupFailedReason           = ClusterReason("CreateSecurityGroupFailedReason")
	CreateSSHKeyFailedReason                  = ClusterReason("CreateSSHKeyFailedReason")
	CreateVmImageFailedReason                 = ClusterReason("CreateVmImageFailedReason")
	CreateVmSpecFailedReason                  = ClusterReason("CreateVmSpecFailedReason")
	AddNodeEntityFailedReason                 = ClusterReason("AddNodeEntityFailedReason")
	SetupBoostrapFailedReason                 = ClusterReason("SetupBoostrapFailedReason")
	SetupHaproxyFailedReason                  = ClusterReason("SetupHaproxyFailedReason")
	InitControlPlaneFailedReason              = ClusterReason("InitControlPlaneFailedReason")
	SetupNetworkCNIFailedReason               = ClusterReason("SetupNetworkCNIFailedReason")
	JoinControlPlaneFailedReason              = ClusterReason("JoinControlPlaneFailedReason")
	JoinWorkerFailedReason                    = ClusterReason("JoinWorkerFailedReason")
)

type Cluster struct {
	Model
	Status struct {
		Phase   ClusterPhase  `json:"phase" enums:"Pending,Provisioning,Provisioned,Failed"`
		Reason  ClusterReason `json:"reason"`
		Message string        `json:"message"`
	} `json:"status"`
	MCIS            string `json:"mcis"`
	Namespace       string `json:"namespace"`
	ClusterConfig   string `json:"clusterConfig"`
	CpLeader        string `json:"cpLeader"`
	NetworkCni      string `json:"networkCni" enums:"canal,kilo"`
	Label           string `json:"label"`
	InstallMonAgent string `json:"installMonAgent" example:"no" default:"yes"`
	Description     string `json:"description"`
	CreatedTime     string `json:"createdTime" example:"2022-01-02T12:00:00Z" default:""`
	Nodes           []Node `json:"nodes"`
}

type ClusterList struct {
	ListModel
	namespace string
	Items     []Cluster `json:"items"`
}

func NewCluster(namespace string, name string) *Cluster {
	return &Cluster{
		Model:     Model{Kind: KIND_CLUSTER, Name: name},
		Namespace: namespace,
		Status: struct {
			Phase   ClusterPhase  "json:\"phase\" enums:\"Pending,Provisioning,Provisioned,Failed\""
			Reason  ClusterReason "json:\"reason\""
			Message string        "json:\"message\""
		}{Phase: ClusterPhasePending, Reason: "", Message: ""},
		Nodes: []Node{},
	}
}

func NewClusterList(namespace string) *ClusterList {
	return &ClusterList{
		ListModel: ListModel{Kind: KIND_CLUSTER_LIST},
		namespace: namespace,
		Items:     []Cluster{},
	}
}

func (self *Cluster) UpdatePhase(phase ClusterPhase) error {
	self.Status.Phase = phase
	if phase != ClusterPhaseFailed {
		self.Status.Reason = ""
		self.Status.Message = ""
	}
	if phase == ClusterPhaseProvisioned {
		self.CreatedTime = lang.GetNowUTC()
	}
	return self.putStore()
}

func (self *Cluster) FailReason(reason ClusterReason, message string) error {
	self.Status.Phase = ClusterPhaseFailed
	self.Status.Reason = reason
	self.Status.Message = message
	return self.putStore()
}

func (self *Cluster) putStore() error {
	key := lang.GetStoreClusterKey(self.Namespace, self.Name)
	value, _ := json.Marshal(self)
	err := common.CBStore.Put(key, string(value))
	if err != nil {
		return err
	}
	return nil
}

func (self *Cluster) Select() (bool, error) {
	exists := false

	key := lang.GetStoreClusterKey(self.Namespace, self.Name)
	keyValue, err := common.CBStore.Get(key)
	if err != nil {
		logger.Errorf("cluster could not be found in the metadata (key=%s, cause=%v)", key, err)
		return exists, errors.New(fmt.Sprintf("The cluster could not be found in the metadata. (namespace=%s, cluster=%s)", self.Namespace, self.Name))
	}
	exists = (keyValue != nil)
	if exists {
		json.Unmarshal([]byte(keyValue.Value), &self)
		err = getClusterNodes(self)
		if err != nil {
			logger.Errorf("node could not be found in the metadata. (key=%s, cause=%v)", key, err)
			return exists, errors.New(fmt.Sprintf("The nodes could not be found in the metadata. (namespace=%s, cluster=%s)", self.Namespace, self.Name))

		}
	}

	return exists, nil
}

func (self *Cluster) Delete() error {
	// delete node
	keyValues, err := common.CBStore.GetList(lang.GetStoreNodeKey(self.Namespace, self.Name, ""), true)
	if err != nil {
		logger.Errorf("nodes could not be found in the metadata (keys=%v, cause=%v)", keyValues, err)
		return errors.New(fmt.Sprintf("The nodes could not be found in the metadata. (namespace=%s, cluster=%s)", self.Namespace, self.Name))

	}
	for _, keyValue := range keyValues {
		err = common.CBStore.Delete(keyValue.Key)
		if err != nil {
			logger.Errorf("failed to delete the node in the metadata (key=%v, cause=%v)", keyValue.Key, err)
			return errors.New(fmt.Sprintf("Failed to delete the node in the metadata (namespace=%s, cluster=%s)", self.Namespace, self.Name))
		}
	}

	// delete cluster
	key := lang.GetStoreClusterKey(self.Namespace, self.Name)
	err = common.CBStore.Delete(key)
	if err != nil {
		logger.Errorf("failed to delete the cluster in the metadata. (key=%v, cause=%v)", key, err)
		return errors.New(fmt.Sprintf("Failed to delete the cluster in the metadata. (namespace=%s, cluster=%s)", self.Namespace, self.Name))
	}

	return nil
}

func (self *ClusterList) SelectList() error {
	keyValues, err := common.CBStore.GetList(lang.GetStoreClusterKey(self.namespace, ""), true)
	if err != nil {
		logger.Errorf("clusters could not be found in the metadata (keys=%v, cause=%v)", keyValues, err)
		return errors.New(fmt.Sprintf("The clusters could not be found in the metadata. (namespace=%s)", self.namespace))
	}
	self.Items = []Cluster{}
	for _, keyValue := range keyValues {
		if !strings.Contains(keyValue.Key, "/nodes") {
			cluster := &Cluster{}
			json.Unmarshal([]byte(keyValue.Value), &cluster)

			err = getClusterNodes(cluster)
			if err != nil {
				logger.Errorf("nodes could not be found in the metadata. (namespace=%s, cluster=%s, cause=%v)", self.namespace, cluster, err)
				return errors.New(fmt.Sprintf("The nodes could not be found in the metadata. (namespace=%s, cluster=%s)", self.namespace, cluster))
			}
			self.Items = append(self.Items, *cluster)
		}
	}

	return nil
}

func getClusterNodes(cluster *Cluster) error {
	nodeKeyValues, err := common.CBStore.GetList(lang.GetStoreNodeKey(cluster.Namespace, cluster.Name, ""), true)
	if err != nil {
		return err
	}
	for _, nodeKeyValue := range nodeKeyValues {
		node := &Node{}
		json.Unmarshal([]byte(nodeKeyValue.Value), &node)
		cluster.Nodes = append(cluster.Nodes, *node)
	}

	return nil
}
