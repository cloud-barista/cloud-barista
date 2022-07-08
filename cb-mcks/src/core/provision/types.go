package provision

import (
	"github.com/cloud-barista/cb-mcks/src/core/app"
	"github.com/cloud-barista/cb-mcks/src/core/model"
)

const (
	REMOTE_TARGET_PATH    = "/tmp"
	CNI_CANAL_FILE        = "addons/canal/canal_v3.20.0.yaml"
	CNI_KILO_CRDS_FILE    = "addons/kilo/crds_v0.3.0.yaml"
	CNI_KILO_KUBEADM_FILE = "addons/kilo/kilo-kubeadm-flannel_v0.3.0.yaml"
	CNI_KILO_FLANNEL_FILE = "addons/kilo/kube-flannel_v0.14.0.yaml"
)

type Machine struct {
	Name       string
	PublicIP   string
	PrivateIP  string
	Username   string
	CSP        app.CSP
	Role       app.ROLE
	Region     string
	Zone       string
	Spec       string
	Credential string
}
type ControlPlaneMachine struct {
	*Machine
}
type WorkerNodeMachine struct {
	*Machine
}

type Provisioner struct {
	Cluster              *model.Cluster
	leader               *ControlPlaneMachine
	ControlPlaneMachines map[string]*ControlPlaneMachine
	WorkerNodeMachines   map[string]*WorkerNodeMachine
}
