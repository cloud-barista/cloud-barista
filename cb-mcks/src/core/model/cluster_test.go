package model

import (
	"testing"
)

const ()

func TestClusterCRUD(t *testing.T) {

	namespace := "namespace-1"
	custerName := "cluster-1"

	// insert
	cluster := NewCluster(namespace, custerName)
	cluster.MCIS = "mcis-1"
	err := cluster.UpdatePhase(ClusterPhaseProvisioning)
	if err != nil {
		t.Fatalf("Cluster.Insert error - TestClusaterCRUD :: mcis (cause=%v)", err)
	}

	// verify insert
	cluster = NewCluster(namespace, custerName)
	_, err = cluster.Select()
	if err != nil {
		t.Fatalf("Cluster.Insert error - TestClusaterCRUD :: cluster (cause=%v)", err)
	}
	if cluster.MCIS != "mcis-1" {
		t.Fatalf("Cluster.Insert verify : not equals mcis (%s != %s)", "mcis-1", cluster.MCIS)
	}

	// verify update
	cluster = NewCluster(namespace, custerName)
	_, err = cluster.Select()
	if err != nil {
		t.Fatalf("Cluster.Update error (cause=%v)", err)
	}
	if cluster.MCIS != "mcis-modifed" {
		t.Fatalf("Cluster.Update verify : not equals mcis (%s != %s)", "mcis-modifed", cluster.MCIS)
	}
	if cluster.Status.Phase != ClusterPhasePending {
		t.Fatalf("Cluster.Update verify : phase is not 'STATUS_PROVISIONING' (%s)", cluster.Status.Phase)
	}

	// delete
	cluster = NewCluster(namespace, custerName)
	err = cluster.Delete()
	if err != nil {
		t.Fatalf("Cluster.Delete error (cause=%v)", err)
	}
	// verify delete
	cluster = NewCluster(namespace, custerName)
	_, err = cluster.Select()
	if err == nil {
		t.Fatalf("Cluster.Delete exist data, not deleted")
	}

}

func TestClusterSelectList(t *testing.T) {

	namespace := "namespace-2"
	custerName := "cluster-2"

	// insert
	cluster := NewCluster(namespace, custerName)
	err := cluster.UpdatePhase(ClusterPhaseProvisioning)
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
	err := cluster.UpdatePhase(ClusterPhaseProvisioned)
	if err != nil {
		t.Fatalf("error cluster.Complete() (cause=%v)", err)
	}
	if cluster.Status.Phase != ClusterPhaseProvisioned {
		t.Fatalf("error cluster.Complete() NOT_STATUS_COMPLETED (cause=%v)", err)
	}
}

func TestClusterFail(t *testing.T) {

	cluster := NewCluster("namespace-4", "cluster-4")

	err := cluster.FailReason(CreateMCISFailedReason, "failed create mcis")
	if err != nil {
		t.Fatalf("error cluster.Fail() (cause=%v)", err)
	}
	if cluster.Status.Phase != ClusterPhaseFailed {
		t.Fatalf("error cluster.Fail() NOT_STATUS_FAILED (cause=%v)", err)
	}
}
