/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloud-barista/cb-operator/src/cmd"
	"github.com/cloud-barista/cb-operator/src/common"
)

func errCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func scanAndWriteMode() {

	fmt.Println("")
	fmt.Println("[Options]")
	fmt.Println("1: Docker Compose environment (Requires Docker and Docker Compose)")
	fmt.Println("2: Kubernetes environment (Requires Kubernetes cluster with Helm 3)")
	fmt.Println("")
	fmt.Print("Choose 1 or 2: ")

	var userInput uint8
	fmt.Scanf("%d", &userInput)

	var tempStr string

	switch userInput {
	case 1:
		fmt.Println("[1: Docker Compose environment (Requires Docker and Docker Compose)] selected.")
		tempStr = common.ModeDockerCompose
	case 2:
		fmt.Println("[2: Kubernetes environment (Requires Kubernetes cluster with Helm 3)] selected.")
		tempStr = common.ModeKubernetes
	default:
		fmt.Println("You should choose between 1 and 2.")
		return
	}

	err := ioutil.WriteFile("./CB_OPERATOR_MODE", []byte(tempStr), os.FileMode(0644))
	errCheck(err)

	fmt.Println("")
	fmt.Println("CB_OPERATOR_MODE is set to: " + tempStr)
	fmt.Println("To change CB_OPERATOR_MODE, just delete the CB_OPERATOR_MODE file and re-run the cb-operator.")
}

func readMode() string {
	if _, err := os.Stat("./CB_OPERATOR_MODE"); err == nil {
		// if file exists
		data, err := ioutil.ReadFile("./CB_OPERATOR_MODE")
		errCheck(err)

		common.CBOperatorMode = string(data)
		fmt.Println("CB_OPERATOR_MODE: " + common.CBOperatorMode)

		//if common.CBOperatorMode == common.DockerCompose || common.CBOperatorMode == common.Kubernetes {
		return common.CBOperatorMode
		//}

	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		fmt.Println("CB_OPERATOR_MODE file does not exist.")
		scanAndWriteMode()
		result := readMode()
		return result

	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

		errCheck(err)
		return ""
	}
	//return ""
}

//var CBOperatorMode string

func main() {

	mode := readMode()

	switch mode {
	case common.ModeDockerCompose:
		cmd.Execute()
	case common.ModeKubernetes:
		cmd.Execute()
	default:
		fmt.Println("Invalid CB_OPERATOR_MODE: " + mode)
		fmt.Println("CB_OPERATOR_MODE should be one of these: " + common.ModeDockerCompose + ", " + common.ModeKubernetes)

		//fmt.Println("To change CB_OPERATOR_MODE, just delete the CB_OPERATOR_MODE file and re-run the cb-operator.")
		scanAndWriteMode()
		main()
	}

}
