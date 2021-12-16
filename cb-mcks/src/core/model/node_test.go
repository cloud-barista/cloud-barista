package model

import (
	"testing"
)

func TestNodeCRD(t *testing.T) {

	namespace := "namespace-1"
	custerName := "cluster-1"
	nodeName := "node-1"

	// insert
	node := NewNode(namespace, custerName, nodeName)
	err := node.Insert()
	if err != nil {
		t.Fatalf("Node.Insert error - node.Insert() (cause=%v)", err)
	}

	// verify insert
	node = NewNode(namespace, custerName, nodeName)
	err = node.Select()
	if err != nil {
		t.Fatalf("Node.Insert error - node.Select() (cause=%v)", err)
	}

	// delete
	node = NewNode(namespace, custerName, nodeName)
	err = node.Delete()
	if err != nil {
		t.Fatalf("Node.Delete error - node.Delete() (cause=%v)", err)
	}
	// verify delete
	node = NewNode(namespace, custerName, nodeName)
	err = node.Select()
	if err == nil {
		t.Fatalf("Node.Delete exist data, not deleted")
	}

}

func TestNodeSelectList(t *testing.T) {

	namespace := "namespace-1"
	custerName := "cluster-1"
	nodeName := "node-1"

	// insert
	node := NewNode(namespace, custerName, nodeName)
	err := node.Insert()
	if err != nil {
		t.Fatalf("Node.Insert error (cause=%v)", err)
	}

	// list
	nodeList := NewNodeList(namespace, custerName)
	err = nodeList.SelectList()
	if err != nil {
		t.Fatalf("error (cause=%v)", err)
	}

	if len(nodeList.Items) != 1 {
		t.Fatalf("missmatched rows (count=%v)", len(nodeList.Items))
	}
	t.Log(len(nodeList.Items))
	t.Log(nodeList.namespace)

}
