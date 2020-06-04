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
	"github.com/spf13/cobra"
	"github.com/cloud-barista/cb-operator/src/common"
)



// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Setup and Run Cloud-Barista System",
	Long: `Setup and Run Cloud-Barista System`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n[Setup and Run Cloud-Barista]\n")

		if common.FileStr == "" {
			fmt.Println("file is required")
		} else {
			/*
			var configuration mcisReq

    		viper.SetConfigFile(fileStr)
			if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config file, %s", err)
			}
			err := viper.Unmarshal(&configuration)
			if err != nil {
			fmt.Printf("Unable to decode into struct, %v", err)
			}

			common.PrintJsonPretty(configuration)
			*/

			cmdStr := "sudo docker-compose -f " + common.FileStr + " up"
			//fmt.Println(cmdStr)
			common.SysCall(cmdStr)
		}

	},
}


func init() {
	rootCmd.AddCommand(runCmd)

	pf := runCmd.PersistentFlags()
	pf.StringVarP(&common.FileStr, "file", "f", "../docker-compose.yaml", "Path to Cloud-Barista Docker-compose file")
//	cobra.MarkFlagRequired(pf, "file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
