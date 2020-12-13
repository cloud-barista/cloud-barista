package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cloud-barista/cb-ladybug/src/utils/config"
	"github.com/cloud-barista/cb-ladybug/src/utils/lang"

	ssh "github.com/cloud-barista/cb-spider/cloud-control-manager/vm-ssh"
)

type VM struct {
	Name         string   `json:"name"`
	Config       string   `json:"connectionName"`
	VPC          string   `json:"vNetId"`
	Subnet       string   `json:"subnetId"`
	Firewall     []string `json:"securityGroupIds"`
	SSHKey       string   `json:"sshKeyId"`
	Image        string   `json:"imageId"`
	Spec         string   `json:"specId"`
	UserAccount  string   `json:"vmUserAccount"`
	UserPassword string   `json:"vmUserPassword"`
	Description  string   `json:"description"`
	PublicIP     string   `json:"publicIP"` // output
	Credential   string   // private
	UId          string   `json:"uid"`
	Role         string   `json:"role"`
}

func (v *VM) ConnectionTest(sshInfo *ssh.SSHInfo, vm *VM) error {
	cmd := "/bin/hostname"
	_, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		fmt.Println("Error while running cmd: "+cmd, err)
		return err
	}
	return nil
}

func (v *VM) CopyScripts(sshInfo *ssh.SSHInfo, vm *VM) error {
	sourcePath := fmt.Sprintf("%s/src/scripts/", config.Config.AppRootPath)
	sourceFile := []string{config.BOOTSTRAP_FILE}
	if vm.Role == config.CONTROL_PLANE {
		sourceFile = append(sourceFile, config.INIT_FILE)
	}
	targetPath := config.Config.TargetPath + "/"

	fmt.Printf("start script file copy (src=%s, dest=%s)\n", sourcePath, targetPath)
	for _, f := range sourceFile {
		src := fmt.Sprintf("%s%s", sourcePath, f)
		dest := fmt.Sprintf("%s%s", targetPath, f)
		if err := ssh.SSHCopy(*sshInfo, src, dest); err != nil {
			fmt.Println("Error while copying file ", err)
			return errors.New("copy scripts error")
		}
	}
	return nil
}

func (v *VM) Bootstrap(sshInfo *ssh.SSHInfo) (bool, error) {
	cmd := fmt.Sprintf("cd %s;./%s", config.Config.TargetPath, config.BOOTSTRAP_FILE)

	result, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		fmt.Println("Error while running cmd: "+cmd, err)
		return false, errors.New("k8s bootstrap error")
	}
	if strings.Contains(result, "kubectl set on hold") {
		return true, nil
	} else {
		return false, nil
	}
}

func (v *VM) ControlPlaneInit(sshInfo *ssh.SSHInfo, ip string) (string, string, error) {
	var workerJoinCmd string

	cmd := fmt.Sprintf("cd %s;./%s", config.Config.TargetPath, config.INIT_FILE)
	cpInitResult, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		fmt.Println("Error while running cmd: "+cmd, err)
		return "", "", errors.New("k8s control plane node init error")
	}
	if strings.Contains(cpInitResult, "Your Kubernetes control-plane has initialized successfully") {
		workerJoinCmd = lang.GetWorkerJoinCmd(cpInitResult)
	} else {
		return "", "", nil
	}

	cmd = fmt.Sprintf("sudo sed '5s/.*/    server: https:\\/\\/%s:6443/g' /etc/kubernetes/admin.conf", ip)
	clusterConfig, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		fmt.Println("Error while running cmd: "+cmd, err)
	}

	return workerJoinCmd, clusterConfig, nil
}

func (v *VM) WorkerJoin(sshInfo *ssh.SSHInfo, workerJoinCmd *string) (bool, error) {
	if *workerJoinCmd == "" {
		return false, errors.New("worker node join command empty")
	}
	cmd := *workerJoinCmd
	result, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		fmt.Println("Error while running cmd: "+cmd, err)
		return false, errors.New("k8s worker node join error")
	}
	if strings.Contains(result, "This node has joined the cluster") {
		return true, nil
	} else {
		return false, errors.New("worker node join failed")
	}
}

func (v *VM) WorkerJoinForAddNode(sshInfo *ssh.SSHInfo, workerJoinCmd *string) (bool, error) {
	if *workerJoinCmd == "" {
		return false, errors.New("worker node join command empty")
	}
	cmd := fmt.Sprintf("sudo %s", *workerJoinCmd)
	result, err := ssh.SSHRun(*sshInfo, cmd)
	if err != nil {
		fmt.Println("Error while running cmd: "+cmd, err)
		return false, errors.New("k8s worker node join error")
	}
	if strings.Contains(result, "This node has joined the cluster") {
		return true, nil
	} else {
		return false, errors.New("worker node join failed")
	}
}
