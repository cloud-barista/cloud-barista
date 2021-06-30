package common

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// FileStr is a variable that holds path to the docker-compose.yaml.
var FileStr string

//var CommandStr string
//var TargetStr string

// CBOperatorMode is a variable that holds current cb-operator's mode.
var CBOperatorMode string

// ModeDockerCompose is a variable that holds string indicating Docker Compose mode.
var ModeDockerCompose string = "DockerCompose"

// ModeKubernetes is a variable that holds string indicating Kubernetes mode.
var ModeKubernetes string = "Kubernetes"

// DefaultDockerComposeConfig is a variable that holds path to docker-compose.yaml
var DefaultDockerComposeConfig string = "../docker-compose-mode-files/docker-compose.yaml"

// DefaultKubernetesConfig is a variable that holds path to helm-chart/values.yaml
var DefaultKubernetesConfig string = "../helm-chart/values.yaml"

// NotDefined is a variable that holds the string "Not_Defined"
var NotDefined string = "Not_Defined"

// CBK8sNamespace is a variable that holds the K8s namespace that CB-Operator will use.
var CBK8sNamespace string = "cloud-barista"

// CBHelmReleaseName is a variable that holds the K8s Helm release name that CB-Operator will use.
var CBHelmReleaseName string = "cloud-barista"

// SysCall executes user-passed command via system call.
func SysCall(cmdStr string) {
	//cmdStr := "sudo docker-compose -f " + common.FileStr + " up"
	cmd := exec.Command("/bin/sh", "-c", cmdStr)

	cmdReader, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// SysCallDockerComposePs executes `docker-compose ps` command via system call.
func SysCallDockerComposePs() {
	fmt.Println("\n[v]Status of Cloud-Barista runtimes")
	cmdStr := "sudo COMPOSE_PROJECT_NAME=cloud-barista docker-compose -f " + FileStr + " ps"
	SysCall(cmdStr)
}

// GenConfigPath receives path-to-config-file and cb-operator-mode, and returns path-to-config-file.
func GenConfigPath(fileStr string, mode string) string {
	returnStr := fileStr
	switch mode {
	case ModeDockerCompose:
		if fileStr == NotDefined {
			returnStr = DefaultDockerComposeConfig
		}
	case ModeKubernetes:
		if fileStr == NotDefined {
			returnStr = DefaultKubernetesConfig
		}
	default:

	}
	fmt.Println("[Config path] " + returnStr)
	fmt.Println()
	return returnStr
}
