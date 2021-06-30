package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cloud-barista/cb-ladybug/src/core/common"
	"github.com/cloud-barista/cb-ladybug/src/utils/lang"
)

const (
	STATUS_CREATED      = "created"
	STATUS_PROVISIONING = "provisioning"
	STATUS_COMPLETED    = "completed"
	STATUS_FAILED       = "failed"
)

type Cluster struct {
	Model
	Status        string `json:"status"`
	UId           string `json:"uid"`
	MCIS          string `json:"mcis"`
	Namespace     string `json:"namespace"`
	ClusterConfig string `json:"clusterConfig"`
	CpLeader      string `json:"cpLeader"`
	NetworkCni    string `json:"networkCni"`
	Nodes         []Node `json:"nodes"`
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
		Nodes:     []Node{},
	}
}

func NewClusterList(namespace string) *ClusterList {
	return &ClusterList{
		ListModel: ListModel{Kind: KIND_CLUSTER_LIST},
		namespace: namespace,
		Items:     []Cluster{},
	}
}

func (self *Cluster) Insert() error {
	self.Status = STATUS_CREATED
	return self.putStore()
}

func (self *Cluster) Update() error {
	self.Status = STATUS_PROVISIONING
	return self.putStore()
}

func (self *Cluster) Complete() error {
	self.Status = STATUS_COMPLETED
	return self.putStore()
}

func (self *Cluster) Fail() error {
	self.Status = STATUS_FAILED
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

func (self *Cluster) Select() error {
	key := lang.GetStoreClusterKey(self.Namespace, self.Name)
	keyValue, err := common.CBStore.Get(key)
	if err != nil {
		return err
	}
	if keyValue == nil {
		return errors.New(fmt.Sprintf("%s not found", key))
	}
	json.Unmarshal([]byte(keyValue.Value), &self)

	err = getClusterNodes(self)
	if err != nil {
		return err
	}

	return nil
}

func (self *Cluster) Delete() error {
	// delete node
	keyValues, err := common.CBStore.GetList(lang.GetStoreNodeKey(self.Namespace, self.Name, ""), true)
	if err != nil {
		return err
	}
	for _, keyValue := range keyValues {
		err = common.CBStore.Delete(keyValue.Key)
		if err != nil {
			return err
		}
	}

	// delete cluster
	key := lang.GetStoreClusterKey(self.Namespace, self.Name)
	err = common.CBStore.Delete(key)
	if err != nil {
		return err
	}

	return nil
}

func (self *ClusterList) SelectList() error {
	keyValues, err := common.CBStore.GetList(lang.GetStoreClusterKey(self.namespace, ""), true)
	if err != nil {
		return err
	}
	self.Items = []Cluster{}
	for _, keyValue := range keyValues {
		if !strings.Contains(keyValue.Key, "/nodes") {
			cluster := &Cluster{}
			json.Unmarshal([]byte(keyValue.Value), &cluster)

			err = getClusterNodes(cluster)
			if err != nil {
				return err
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
