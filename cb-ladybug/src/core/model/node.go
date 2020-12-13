package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cloud-barista/cb-ladybug/src/core/common"
	"github.com/cloud-barista/cb-ladybug/src/utils/lang"
)

type Node struct {
	Model
	Credential string `json:"credential"`
	PublicIP   string `json:"publicIp"`
	UId        string `json:"uid"`
	Role       string `json:"role"`
	Spec       string `json:"spec"`
}

type NodeList struct {
	Kind  string `json:"kind"`
	Items []Node `json:"items"`
}

type NodeReq struct {
	Config          string `json:"config"`
	WorkerNodeSpec  string `json:"workerNodeSpec"`
	WorkerNodeCount int    `json:"workerNodeCount"`
}

func NewNode(vm VM) *Node {
	return &Node{
		Model:      Model{Kind: KIND_NODE, Name: vm.Name},
		Credential: vm.Credential,
		PublicIP:   vm.PublicIP,
		UId:        vm.UId,
		Role:       vm.Role,
		Spec:       vm.Spec,
	}
}

func NewNodeDef(nodeName string) *Node {
	return &Node{
		Model: Model{Kind: KIND_NODE, Name: nodeName},
	}
}

func NewNodeList() *NodeList {
	return &NodeList{
		Kind:  KIND_NODE_LIST,
		Items: []Node{},
	}
}

func (n *Node) Select(namespace string, clusterName string, node *Node) (*Node, error) {
	key := lang.GetStoreNodeKey(namespace, clusterName, node.Name)
	keyValue, err := common.CBStore.Get(key)
	if keyValue == nil {
		return nil, errors.New(fmt.Sprintf("%s not found", key))
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(keyValue.Value), &node)
	return node, nil
}

func (n *Node) Insert(namespace string, clusterName string, node *Node) (*Node, error) {
	key := lang.GetStoreNodeKey(namespace, clusterName, node.Name)
	value, _ := json.Marshal(node)
	err := common.CBStore.Put(key, string(value))
	if err != nil {
		common.CBLog.Error(err)
	}

	return node, nil
}

func (n *Node) Delete(namespace string, clusterName string, node *Node) error {
	key := lang.GetStoreNodeKey(namespace, clusterName, node.Name)
	err := common.CBStore.Delete(key)
	if err != nil {
		return err
	}

	return nil
}

func (n *NodeList) SelectList(namespace string, clusterName string, nodes *NodeList) (*NodeList, error) {
	keyValues, err := common.CBStore.GetList(lang.GetStoreNodeKey(namespace, clusterName, ""), true)
	if err != nil {
		return nil, err
	}
	for _, keyValue := range keyValues {
		node := &Node{}
		json.Unmarshal([]byte(keyValue.Value), &node)
		nodes.Items = append(nodes.Items, *node)
	}

	return nodes, nil
}
