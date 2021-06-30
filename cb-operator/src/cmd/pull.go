package cmd

import (
	"fmt"

	"github.com/cloud-barista/cb-operator/src/common"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull images of Cloud-Barista System containers",
	Long:  `Pull images of Cloud-Barista System containers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n[Pull images of Cloud-Barista System containers]")
		fmt.Println()

		if common.FileStr == "" {
			fmt.Println("file is required")
		} else {
			common.FileStr = common.GenConfigPath(common.FileStr, common.CBOperatorMode)
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

			cmdStr := "sudo COMPOSE_PROJECT_NAME=cloud-barista docker-compose -f " + common.FileStr + " pull"
			//fmt.Println(cmdStr)
			common.SysCall(cmdStr)
		}

	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pf := pullCmd.PersistentFlags()
	pf.StringVarP(&common.FileStr, "file", "f", common.NotDefined, "User-defined configuration file")
	//	cobra.MarkFlagRequired(pf, "file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
