package common

import (
	"bufio"
	"fmt"
	"os/exec"
)

var FileStr string
var CommandStr string
var TargetStr string
var CB_OPERATOR_MODE string

const (
	Mode_DockerCompose           string = "DockerCompose"
	Mode_Kubernetes              string = "Kubernetes"
	Default_DockerCompose_Config string = "../docker-compose-mode-files/docker-compose.yaml"
	Default_Kubernetes_Config    string = "../helm-chart/values.yaml"
	Not_Defined                  string = "Not_Defined"

	CB_K8s_Namespace     string = "cloud-barista"
	CB_Helm_Release_Name string = "cloud-barista"
)

func SysCall(cmdStr string) {
	//cmdStr := "sudo docker-compose -f " + common.FileStr + " up"
	cmd := exec.Command("/bin/sh", "-c", cmdStr)

	cmdReader, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	cmd.Start()
	scanner := bufio.NewScanner(cmdReader)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
	cmd.Wait()
}

func SysCall_docker_compose_ps() {
	fmt.Println("\n[v]Status of Cloud-Barista runtimes")
	cmdStr := "sudo COMPOSE_PROJECT_NAME=cloud-barista docker-compose -f " + FileStr + " ps"
	SysCall(cmdStr)
}

func GenConfigPath(fileStr string, mode string) string {
	returnStr := fileStr
	switch mode {
	case Mode_DockerCompose:
		if fileStr == Not_Defined {
			returnStr = Default_DockerCompose_Config
		}
	case Mode_Kubernetes:
		if fileStr == Not_Defined {
			returnStr = Default_Kubernetes_Config
		}
	default:

	}
	fmt.Println("[Config path] " + returnStr + "\n")
	return returnStr
}
