package common

import (
	uuid "github.com/google/uuid"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)


var FileStr string
var CommandStr string
var TargetStr string

func SysCall(cmdStr string) (string, error) {
	cmd := exec.Command("bash", "-c", cmdStr)
	
	cmdOut, err := cmd.Output()

	if err != nil {
		//panic(err)
		return "nil", err
	}
	fmt.Println(string(cmdOut))

	return string(cmdOut), nil
	
}

func SysLookPath(cmdStr string) (string, error) {
	path, err := exec.LookPath(cmdStr)

	if err != nil {
		//panic(err)
		return "nil", err
	}
	fmt.Println(string(path))

	return string(path), nil
	
}

func GenUuid() string {
	return uuid.New().String()
}

func GetNumCPU() int {
	return runtime.NumCPU()
}


func PrintJsonPretty(v interface{}) {
	prettyJSON, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		fmt.Printf("%+v\n", v)
	} else {
		fmt.Printf("%s\n", string(prettyJSON))
	}
}
