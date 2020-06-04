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

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Stop and Remove Cloud-Barista System",
	Long: `Stop and Remove Cloud-Barista System. Stop and Remove Cloud-Barista runtimes and related container images and meta-DB if necessary`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("\n[Remove Cloud-Barista]\n")

		if common.FileStr == "" {
			fmt.Println("file is required")
		} else {

			cmdStr := ""
			if volFlag && imgFlag {
				cmdStr = "sudo docker-compose -f " + common.FileStr + " down -v --rmi all"
			} else if volFlag {
				cmdStr = "sudo docker-compose -f " + common.FileStr + " down -v"
			} else if imgFlag {
				cmdStr = "sudo docker-compose -f " + common.FileStr + " down --rmi all"
			} else {
				cmdStr = "sudo docker-compose -f " + common.FileStr + " down"
			}

			//fmt.Println(cmdStr)
			common.SysCall(cmdStr)

			fmt.Println("\n[v]Status of Cloud-Barista runtimes")
			cmdStr = "sudo docker-compose ps"
			//fmt.Println(cmdStr)
			common.SysCall(cmdStr)

		}


	},
}

var volFlag bool
var imgFlag bool

func init() {
	rootCmd.AddCommand(removeCmd)

	pf := removeCmd.PersistentFlags()
	pf.StringVarP(&common.FileStr, "file", "f", "../docker-compose.yaml", "Path to Cloud-Barista Docker-compose file")
//	cobra.MarkFlagRequired(pf, "file")

	pf.BoolVarP(&volFlag, "volumes", "v", false, "Remove named volumes declared in the volumes section of the Compose file")
	pf.BoolVarP(&imgFlag, "images", "i", false, "Remove all images")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
