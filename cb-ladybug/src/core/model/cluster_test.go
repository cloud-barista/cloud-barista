package model

import (
	"testing"

	"github.com/cloud-barista/cb-ladybug/src/utils/lang"
)

const ()

func TestClusterCRUD(t *testing.T) {

	namespace := "namespace-1"
	custerName := "cluster-1"

	uid := lang.GetUid()

	// insert
	cluster := NewCluster(namespace, custerName)
	cluster.UId = uid
	cluster.MCIS = "mcis-1"
	err := cluster.Insert()
	if err != nil {
		t.Fatalf("Cluster.Insert error - TestClusaterCRUD :: mcis (cause=%v)", err)
	}

	// verify insert
	cluster = NewCluster(namespace, custerName)
	err = cluster.Select()
	if err != nil {
		t.Fatalf("Cluster.Insert error - TestClusaterCRUD :: cluster (cause=%v)", err)
	}
	if cluster.UId != uid {
		t.Fatalf("Cluster.Insert verify : not equals uid (%s != %s)", uid, cluster.UId)
	}
	if cluster.Status != STATUS_CREATED {
		t.Fatalf("Cluster.Insert verify : status is not 'STATUS_CREATED' (%s)", cluster.Status)
	}

	// update
	cluster = NewCluster(namespace, custerName)
	cluster.MCIS = "mcis-modifed"
	err = cluster.Update()
	if err != nil {
		t.Fatalf("Cluster.Update error (cause=%v)", err)
	}

	// verify update
	cluster = NewCluster(namespace, custerName)
	err = cluster.Select()
	if err != nil {
		t.Fatalf("Cluster.Update error (cause=%v)", err)
	}
	if cluster.MCIS != "mcis-modifed" {
		t.Fatalf("Cluster.Update verify : not equals uid (%s != %s)", uid, cluster.UId)
	}
	if cluster.Status != STATUS_PROVISIONING {
		t.Fatalf("Cluster.Update verify : status is not 'STATUS_PROVISIONING' (%s)", cluster.Status)
	}

	// delete
	cluster = NewCluster(namespace, custerName)
	err = cluster.Delete()
	if err != nil {
		t.Fatalf("Cluster.Delete error (cause=%v)", err)
	}
	// verify delete
	cluster = NewCluster(namespace, custerName)
	err = cluster.Select()
	if err == nil {
		t.Fatalf("Cluster.Delete exist data, not deleted")
	}

}

func TestClusterSelectList(t *testing.T) {

	namespace := "namespace-2"
	custerName := "cluster-2"

	// insert
	cluster := NewCluster(namespace, custerName)
	err := cluster.Insert()
	if err != nil {
		t.Fatalf("Cluster.Insert error - TestClusterSelectList :: cluster-2 (cause=%v)", err)
	}

	clusterList := NewClusterList(namespace)
	err = clusterList.SelectList()
	if err != nil {
		t.Fatalf("error clusterList.SelectList() (cause=%v)", err)
	}

	if len(clusterList.Items) != 1 {
		t.Fatalf("missmatched rows (count=%v)", len(clusterList.Items))
	}
	t.Log(len(clusterList.Items))
	t.Log(clusterList.namespace)

}

func TestClusterComplete(t *testing.T) {

	cluster := NewCluster("namespace-3", "cluster-3")
	err := cluster.Complete()
	if err != nil {
		t.Fatalf("error cluster.Complete() (cause=%v)", err)
	}
	if cluster.Status != STATUS_COMPLETED {
		t.Fatalf("error cluster.Complete() NOT_STATUS_COMPLETED (cause=%v)", err)
	}
}

func TestClusterFail(t *testing.T) {

	cluster := NewCluster("namespace-4", "cluster-4")

	err := cluster.Fail()
	if err != nil {
		t.Fatalf("error cluster.Fail() (cause=%v)", err)
	}
	if cluster.Status != STATUS_FAILED {
		t.Fatalf("error cluster.Fail() NOT_STATUS_FAILED (cause=%v)", err)
	}
}
