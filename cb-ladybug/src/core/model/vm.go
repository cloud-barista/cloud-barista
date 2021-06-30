package model

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	ssh "github.com/cloud-barista/cb-spider/cloud-control-manager/vm-ssh"

	logger "github.com/sirupsen/logrus"
)

type VM struct {
	Name         string     `json:"name"`
	Config       string     `json:"connectionName"`
	VPC          string     `json:"vNetId"`
	Subnet       string     `json:"subnetId"`
	Firewall     []string   `json:"securityGroupIds"`
	SSHKey       string     `json:"sshKeyId"`
	Image        string     `json:"imageId"`
	Spec         string     `json:"specId"`
	UserAccount  string     `json:"vmUserAccount"`
	UserPassword string     `json:"vmUserPassword"`
	Description  string     `json:"description"`
	PublicIP     string     `json:"publicIP"`  // output
	PrivateIP    string     `json:"privateIP"` // output
	Credential   string     // private
	Role         string     `json:"role"`
	Csp          config.CSP `json:"csp"`
	IsCPLeader   bool       `json:"isCPLeader"`
}

type VMInfo struct {
	Name       string     `json:"name"`
	Credential string     // private
	Role       string     `json:"role"`
	Csp        config.CSP `json:"csp"`
	IsCPLeader bool       `json:"isCPLeader"`
}

const (
	remoteTargetPath = "/tmp"
)

func (self *VM) ConnectionTest(sshInfo *ssh.SSHInfo) error {
	cmd := "/bin/hostname"
	_, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		logger.Warnf("connection test error (server=%s, cause=%s)", sshInfo.ServerPort, err)
		return err
	}
	return nil
}

func (self *VM) CopyScripts(sshInfo *ssh.SSHInfo, networkCni string) error {
	sourcePath := fmt.Sprintf("%s/src/scripts", *config.Config.AppRootPath)
	sourceFile := []string{config.BOOTSTRAP_FILE}
	if self.Role == config.CONTROL_PLANE && self.IsCPLeader {
		sourceFile = append(sourceFile, config.INIT_FILE)
		sourceFile = append(sourceFile, config.HA_PROXY_FILE)
	}
	if networkCni == config.NETWORKCNI_CANAL {
		sourceFile = append(sourceFile, config.LADYBUG_BOOTSTRAP_CANAL_FILE)
	} else {
		sourceFile = append(sourceFile, config.LADYBUG_BOOTSTRAP_KILO_FILE)
	}
	sourceFile = append(sourceFile, config.SYSTEMD_SERVICE_FILE)

	logger.Infof("start script file copy (vm=%s, src=%s, dest=%s)\n", self.Name, sourcePath, remoteTargetPath)
	for _, f := range sourceFile {
		src := fmt.Sprintf("%s/%s", sourcePath, f)
		dest := fmt.Sprintf("%s/%s", remoteTargetPath, f)
		if err := ssh.SSHCopy(*sshInfo, src, dest); err != nil {
			return errors.New(fmt.Sprintf("copy scripts error (server=%s, cause=%s)", sshInfo.ServerPort, err))
		}
	}
	logger.Infof("end script file copy (vm=%s, server=%s)\n", self.Name, sshInfo.ServerPort)
	return nil
}

func (self *VM) SetSystemd(sshInfo *ssh.SSHInfo, networkCni string) error {
	var bsFile string
	if networkCni == config.NETWORKCNI_CANAL {
		bsFile = config.LADYBUG_BOOTSTRAP_CANAL_FILE
	} else {
		bsFile = config.LADYBUG_BOOTSTRAP_KILO_FILE
	}

	cmd := fmt.Sprintf("cd %s;./%s", remoteTargetPath, bsFile)
	_, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		return errors.New(fmt.Sprintf("create ladybug-bootstrap error (name=%s)", self.Name))
	}

	cmd = fmt.Sprintf("cd %s;./%s", remoteTargetPath, config.SYSTEMD_SERVICE_FILE)
	_, err = ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		return errors.New(fmt.Sprintf("set systemd service error (name=%s)", self.Name))
	}
	return nil
}

func (self *VM) Bootstrap(sshInfo *ssh.SSHInfo) error {
	cmd := fmt.Sprintf("cd %s;./%s", remoteTargetPath, config.BOOTSTRAP_FILE)

	result, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		return errors.New("k8s bootstrap error")
	}
	if strings.Contains(result, "kubectl set on hold") {
		return nil
	} else {
		return errors.New(fmt.Sprintf("k8s bootstrap failed (name=%s)", self.Name))
	}
}

func (self *VM) InstallHAProxy(sshInfo *ssh.SSHInfo, IPs []string) error {
	var servers string
	for i, ip := range IPs {
		servers += fmt.Sprintf("  server  api%d  %s:6443  check", i+1, ip)
		if i < len(IPs)-1 {
			servers += "\\n"
		}
	}
	cmd := fmt.Sprintf("sudo sed 's/^{{SERVERS}}/%s/g' %s/%s", servers, remoteTargetPath, config.HA_PROXY_FILE)
	result, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		logger.Warnf("get haproxy command error (name=%s, cause=%v)", self.Name, err)
		return err
	}
	_, err = ssh.SSHRun(*sshInfo, result)
	if err != nil {
		logger.Warnf("install haproxy error (name=%s, cause=%v)", self.Name, err)
		return err
	}
	return nil
}

func (self *VM) ControlPlaneInit(sshInfo *ssh.SSHInfo, reqKubernetes Kubernetes) ([]string, string, error) {
	var joinCmd []string

	cmd := fmt.Sprintf("cd %s;./%s %s %s %s", remoteTargetPath, config.INIT_FILE, reqKubernetes.PodCidr, reqKubernetes.ServiceCidr, reqKubernetes.ServiceDnsDomain)
	cpInitResult, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		logger.Warnf("control plane init error (name=%s, cause=%v)", self.Name, err)
		return nil, "", errors.New("k8s control plane node init error")
	}
	if strings.Contains(cpInitResult, "Your Kubernetes control-plane has initialized successfully") {
		joinCmd = getJoinCmd(cpInitResult)
	} else {
		return nil, "", errors.New(fmt.Sprintf("control palne init failed (name=%s)", self.Name))
	}

	cmd = "sudo cat /etc/kubernetes/admin.conf"
	clusterConfig, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		logger.Errorf("Error while running cmd %s (vm=%s, cause=%v)", cmd, self.Name, err)
	}

	return joinCmd, clusterConfig, nil
}

func (self *VM) InstallNetworkCNI(sshInfo *ssh.SSHInfo, networkCni string) error {
	var cmd string
	if networkCni == config.NETWORKCNI_CANAL {
		cmd = "sudo kubectl apply -f https://docs.projectcalico.org/manifests/canal.yaml --kubeconfig=/etc/kubernetes/admin.conf"
	} else {
		cmd = `sudo kubectl apply -f https://raw.githubusercontent.com/squat/kilo/main/manifests/crds.yaml --kubeconfig=/etc/kubernetes/admin.conf;
		sudo kubectl apply -f https://raw.githubusercontent.com/squat/kilo/master/manifests/kilo-kubeadm-flannel.yaml --kubeconfig=/etc/kubernetes/admin.conf;
		sudo kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml --kubeconfig=/etc/kubernetes/admin.conf;`
	}

	_, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		logger.Warnf("networkCNI install failed (name=%s, cause=%v)", self.Name, err)
		return errors.New("NetworkCNI Install error")
	}
	return nil
}

func (self *VM) ControlPlaneJoin(sshInfo *ssh.SSHInfo, CPJoinCmd *string) error {
	if *CPJoinCmd == "" {
		return errors.New("control-plane node join command empty")
	}
	cmd := fmt.Sprintf("sudo %s", *CPJoinCmd)
	result, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		logger.Warnf("control-plane join error (name=%s, cause=%v)", self.Name, err)
		return errors.New("control-plane node join error")
	}

	if strings.Contains(result, "This node has joined the cluster") {
		_, err = ssh.SSHRun(*sshInfo, "sudo systemctl restart ladybug-bootstrap")
		if err != nil {
			logger.Warnf("ladybug-bootstrap restart error (name=%s, cause=%v)", self.Name, err)
		}
		return nil
	} else {
		logger.Warnf("control-plane join failed (name=%s)", self.Name)
		return errors.New(fmt.Sprintf("control-plane join failed (name=%s)", self.Name))
	}
}

func (self *VM) WorkerJoin(sshInfo *ssh.SSHInfo, workerJoinCmd *string) error {
	if *workerJoinCmd == "" {
		return errors.New("worker node join command empty")
	}
	cmd := fmt.Sprintf("sudo %s", *workerJoinCmd)
	result, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		logger.Warnf("worker join error (name=%s, cause=%v)", self.Name, err)
		return errors.New(fmt.Sprintf("worker node join error (name=%s)", self.Name))
	}
	if strings.Contains(result, "This node has joined the cluster") {
		_, err = ssh.SSHRun(*sshInfo, "sudo systemctl restart ladybug-bootstrap")
		if err != nil {
			logger.Warnf("ladybug-bootstrap restart error (name=%s, cause=%v)", self.Name, err)
		}
		return nil
	} else {
		logger.Warnf("worker join failed (name=%s)", self.Name)
		return errors.New(fmt.Sprintf("worker node join failed (name=%s)", self.Name))
	}
}

func getJoinCmd(cpInitResult string) []string {
	var join1, join2, join3 string
	joinRegex, _ := regexp.Compile("kubeadm\\sjoin\\s(.*?)\\s--token\\s(.*?)\\n")
	joinRegex2, _ := regexp.Compile("--discovery-token-ca-cert-hash\\ssha256:(.*?)\\n")
	joinRegex3, _ := regexp.Compile("--control-plane --certificate-key(.*?)\\n")

	if joinRegex.MatchString(cpInitResult) {
		join1 = joinRegex.FindString(cpInitResult)
	}
	if joinRegex2.MatchString(cpInitResult) {
		join2 = joinRegex2.FindString(cpInitResult)
	}
	if joinRegex3.MatchString(cpInitResult) {
		join3 = joinRegex3.FindString(cpInitResult)
	}

	return []string{fmt.Sprintf("%s %s %s", join1, join2, join3), fmt.Sprintf("%s %s", join1, join2)}
}
