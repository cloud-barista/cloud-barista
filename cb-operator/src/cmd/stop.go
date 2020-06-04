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

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Cloud-Barista System",
	Long: `Stop Cloud-Barista System`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n[Stop Cloud-Barista]\n")

		if common.FileStr == "" {
			fmt.Println("file is required")
		} else {
			cmdStr := "sudo docker-compose -f " + common.FileStr + " stop"
			//fmt.Println(cmdStr)
			common.SysCall(cmdStr)

			fmt.Println("\n[v]Status of Cloud-Barista runtimes")
			cmdStr = "sudo docker-compose ps"
			common.SysCall(cmdStr)
		}

	},
}



func init() {
	rootCmd.AddCommand(stopCmd)

	pf := stopCmd.PersistentFlags()
	pf.StringVarP(&common.FileStr, "file", "f", "../docker-compose.yaml", "Path to Cloud-Barista Docker-compose file")
//	cobra.MarkFlagRequired(pf, "file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
