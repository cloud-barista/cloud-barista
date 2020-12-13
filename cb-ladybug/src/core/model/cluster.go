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
	Nodes         []Node `json:"nodes"`
}

type ClusterList struct {
	Kind  string    `json:"kind"`
	Items []Cluster `json:"items"`
}

type ClusterReq struct {
	Name                  string `json:"name"`
	ControlPlaneNodeSpec  string `json:"controlPlaneNodeSpec"`
	ControlPlaneNodeCount int    `json:"controlPlaneNodeCount"`
	WorkerNodeSpec        string `json:"workerNodeSpec"`
	WorkerNodeCount       int    `json:"workerNodeCount"`
}

func NewCluster(namespace string, name string) *Cluster {
	return &Cluster{
		Model:     Model{Kind: KIND_CLUSTER, Name: name},
		Namespace: namespace,
		Nodes:     []Node{},
	}
}

func NewClusterList() *ClusterList {
	return &ClusterList{
		Kind:  KIND_CLUSTER_LIST,
		Items: []Cluster{},
	}
}

func (c *Cluster) Insert(cluster *Cluster) {
	cluster.Status = STATUS_CREATED
	c.putStore(cluster)
}

func (c *Cluster) Update(cluster *Cluster) {
	cluster.Status = STATUS_PROVISIONING
	c.putStore(cluster)
}

func (c *Cluster) Complete(cluster *Cluster) {
	cluster.Status = STATUS_COMPLETED
	c.putStore(cluster)
}

func (c *Cluster) Fail(cluster *Cluster) {
	cluster.Status = STATUS_FAILED
	c.putStore(cluster)
}

func (c *Cluster) putStore(cluster *Cluster) {
	key := lang.GetStoreClusterKey(cluster.Namespace, cluster.Name)
	value, _ := json.Marshal(cluster)
	err := common.CBStore.Put(key, string(value))
	if err != nil {
		common.CBLog.Error(err)
	}
}

func (c *Cluster) Select(cluster *Cluster) (*Cluster, error) {
	key := lang.GetStoreClusterKey(cluster.Namespace, cluster.Name)
	keyValue, err := common.CBStore.Get(key)

	if keyValue == nil {
		return nil, errors.New(fmt.Sprintf("%s not found", key))
	}
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(keyValue.Value), &cluster)

	cluster, err = getClustersNodes(cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (c *Cluster) Delete(cluster *Cluster) error {
	// delete node
	keyValues, err := common.CBStore.GetList(lang.GetStoreNodeKey(cluster.Namespace, cluster.Name, ""), true)
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
	key := lang.GetStoreClusterKey(cluster.Namespace, cluster.Name)
	err = common.CBStore.Delete(key)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClusterList) SelectList(namespace string, clusters *ClusterList) (*ClusterList, error) {
	keyValues, err := common.CBStore.GetList(lang.GetStoreClusterKey(namespace, ""), true)
	if err != nil {
		return nil, err
	}
	for _, keyValue := range keyValues {
		if !strings.Contains(keyValue.Key, "/nodes") {
			cluster := &Cluster{}
			json.Unmarshal([]byte(keyValue.Value), &cluster)

			cluster, err = getClustersNodes(cluster)
			if err != nil {
				return nil, err
			}
			clusters.Items = append(clusters.Items, *cluster)
		}
	}

	return clusters, nil
}

func getClustersNodes(cluster *Cluster) (*Cluster, error) {
	nodeKeyValues, err := common.CBStore.GetList(lang.GetStoreNodeKey(cluster.Namespace, cluster.Name, ""), true)
	if err != nil {
		return nil, err
	}
	for _, nodeKeyValue := range nodeKeyValues {
		node := &Node{}
		json.Unmarshal([]byte(nodeKeyValue.Value), &node)
		cluster.Nodes = append(cluster.Nodes, *node)
	}

	return cluster, nil
}
