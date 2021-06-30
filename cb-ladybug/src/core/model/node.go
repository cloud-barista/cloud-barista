package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cloud-barista/cb-ladybug/src/core/common"
	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	"github.com/cloud-barista/cb-ladybug/src/utils/lang"
)

type Node struct {
	Model
	namespace   string
	clusterName string
	Credential  string     `json:"credential"`
	PublicIP    string     `json:"publicIp"`
	UId         string     `json:"uid"`
	Role        string     `json:"role"`
	Spec        string     `json:"spec"`
	Csp         config.CSP `json:"csp"`
}

type NodeList struct {
	ListModel
	namespace   string
	clusterName string
	Items       []Node `json:"items"`
}

func NewNodeVM(namespace string, clusterName string, vm VM) *Node {
	return &Node{
		Model:       Model{Kind: KIND_NODE, Name: vm.Name},
		Credential:  vm.Credential,
		PublicIP:    vm.PublicIP,
		Role:        vm.Role,
		Spec:        vm.Spec,
		Csp:         vm.Csp,
		namespace:   namespace,
		clusterName: clusterName,
	}
}

func NewNode(namespace string, clusterName string, nodeName string) *Node {
	return &Node{
		Model:       Model{Kind: KIND_NODE, Name: nodeName},
		namespace:   namespace,
		clusterName: clusterName,
	}
}

func NewNodeList(namespace string, clusterName string) *NodeList {
	return &NodeList{
		ListModel:   ListModel{Kind: KIND_NODE_LIST},
		Items:       []Node{},
		namespace:   namespace,
		clusterName: clusterName,
	}
}

func (self *Node) Select() error {
	key := lang.GetStoreNodeKey(self.namespace, self.clusterName, self.Name)
	keyValue, err := common.CBStore.Get(key)
	if err != nil {
		return err
	}
	if keyValue == nil {
		return errors.New(fmt.Sprintf("%s not found", key))
	}

	json.Unmarshal([]byte(keyValue.Value), &self)
	return nil
}

func (self *Node) Insert() error {
	key := lang.GetStoreNodeKey(self.namespace, self.clusterName, self.Name)
	value, _ := json.Marshal(self)
	err := common.CBStore.Put(key, string(value))
	if err != nil {
		return err
	}

	return nil
}

func (self *Node) Delete() error {
	key := lang.GetStoreNodeKey(self.namespace, self.clusterName, self.Name)
	err := common.CBStore.Delete(key)
	if err != nil {
		return err
	}

	return nil
}

func (self *NodeList) SelectList() error {
	keyValues, err := common.CBStore.GetList(lang.GetStoreNodeKey(self.namespace, self.clusterName, ""), true)
	if err != nil {
		return err
	}
	for _, keyValue := range keyValues {
		node := &Node{}
		json.Unmarshal([]byte(keyValue.Value), &node)
		self.Items = append(self.Items, *node)
	}

	return nil
}
