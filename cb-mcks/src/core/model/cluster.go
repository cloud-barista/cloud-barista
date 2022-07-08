package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloud-barista/cb-mcks/src/core/app"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"
)

/* new instance of cluster-entity */
func NewCluster(namespace string, name string) *Cluster {
	return &Cluster{
		Model:     Model{Kind: app.KIND_CLUSTER, Name: name},
		Namespace: namespace,
		Status:    ClusterStatus{Phase: ClusterPhasePending, Reason: "", Message: ""},
		Nodes:     []*Node{},
	}
}

/* new instance of cluster-entity list */
func NewClusterList(namespace string) *ClusterList {
	return &ClusterList{
		ListModel: ListModel{Kind: app.KIND_CLUSTER_LIST},
		namespace: namespace,
		Items:     []Cluster{},
	}
}

/* new instance of node-entity */
func NewNode(namespace string, clusterName string, nodeName string) *Node {
	return &Node{
		Model:       Model{Kind: app.KIND_NODE, Name: nodeName},
		namespace:   namespace,
		clusterName: clusterName,
	}
}

/* new instance of node-entity list */
func NewNodeList(namespace string, clusterName string) *NodeList {
	return &NodeList{
		ListModel:   ListModel{Kind: app.KIND_NODE_LIST},
		Items:       []*Node{},
		namespace:   namespace,
		clusterName: clusterName,
	}
}

/* cluster-entity */
func (self *Cluster) UpdatePhase(phase ClusterPhase) error {
	self.Status.Phase = phase
	if phase != ClusterPhaseFailed {
		self.Status.Reason = ""
		self.Status.Message = ""
	}
	if phase == ClusterPhaseProvisioned {
		self.CreatedTime = lang.GetNowUTC()
	}
	return self.PutStore()
}

func (self *Cluster) FailReason(reason ClusterReason, message string) error {
	self.Status.Phase = ClusterPhaseFailed
	self.Status.Reason = reason
	self.Status.Message = message
	return self.PutStore()
}

func (self *Cluster) PutStore() error {
	key := getStoreClusterKey(self.Namespace, self.Name)
	value, _ := json.Marshal(self)
	err := app.CBStore.Put(key, string(value))
	if err != nil {
		return err
	}
	return nil
}

func (self *Cluster) Select() (bool, error) {
	exists := false

	key := getStoreClusterKey(self.Namespace, self.Name)
	keyValue, err := app.CBStore.Get(key)
	if err != nil {
		return exists, err
	}
	exists = (keyValue != nil)
	if exists {
		json.Unmarshal([]byte(keyValue.Value), &self)
	}

	return exists, nil
}

func (self *Cluster) Delete() error {

	// delete cluster
	key := getStoreClusterKey(self.Namespace, self.Name)
	if err := app.CBStore.Delete(key); err != nil {
		return err
	}

	return nil
}

func (self *Cluster) NextNodeIndex(role app.ROLE) int {

	max := 0
	for _, node := range self.Nodes {
		if node.Role == role {
			idx := lang.GetNodeNameIndex(node.Name)
			if idx > max {
				max = idx
			}
		}
	}
	return (max + 1)
}

func (self *Cluster) GetNode(nodeName string) *Node {

	for _, node := range self.Nodes {
		if node.Name == nodeName {
			return node
		}
	}
	return nil
}

func (self *Cluster) DeleteNode(nodeName string) error {

	for i, node := range self.Nodes {
		if node.Name == nodeName {
			self.Nodes = append(self.Nodes[:i], self.Nodes[i+1:]...)
		}
	}
	if err := self.PutStore(); err != nil {
		return err
	}

	return nil
}
func (self *Cluster) ExistsNode(nodeName string) bool {

	var exists = false
	for _, node := range self.Nodes {
		if node.Name == nodeName {
			exists = true
		}
	}
	return exists

}

func (self *ClusterList) SelectList() error {
	keyValues, err := app.CBStore.GetList(getStoreClusterKey(self.namespace, ""), true)
	if err != nil {
		return err
	}
	self.Items = []Cluster{}
	for _, keyValue := range keyValues {
		if !strings.Contains(keyValue.Key, "/nodes") {
			cluster := &Cluster{}
			json.Unmarshal([]byte(keyValue.Value), &cluster)
			self.Items = append(self.Items, *cluster)
		}
	}

	return nil
}

// get store cluster key
func getStoreClusterKey(namespace string, clusterName string) string {
	if clusterName == "" {
		return fmt.Sprintf("/ns/%s/clusters", namespace)
	} else {
		return fmt.Sprintf("/ns/%s/clusters/%s", namespace, clusterName)
	}
}
