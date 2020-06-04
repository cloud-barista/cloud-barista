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
package cmd

import (
	"fmt"
	"github.com/cloud-barista/cb-operator/src/common"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run commands in a target component of Cloud-Barista System",
	Long: `Run commands in your components of Cloud-Barista System. 
	For instance, you can get an interactive prompt of cb-tumblebug by
	[operator exec cb-tumblebug sh]`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n[Execute COMMAND{"+common.CommandStr+"] in the TARGET{"+common.TargetStr + "}]\n")

		if common.TargetStr != "" && common.CommandStr != "" {	
			//Need to resolve a problem which "sh" command can not get interactive shell. (-T mode is added intentionally)		
			cmdStr := "sudo docker-compose exec -T " + common.TargetStr + " " + common.CommandStr
			fmt.Println(cmdStr)
			common.SysCall(cmdStr)
		} else {
			fmt.Println("Need to provide -t [target name] and -c [command]")
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	common.TargetStr = ""
	common.CommandStr = ""
	execCmd.PersistentFlags().StringVarP(&common.TargetStr, "target", "t", "", "Name of CB component to command")
	execCmd.PersistentFlags().StringVarP(&common.CommandStr, "command", "c", "", "Command to excute")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
