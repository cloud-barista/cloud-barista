package cmd

import (
	"fmt"

	"github.com/cloud-barista/cb-operator/src/common"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Setup and Run Cloud-Barista System",
	Long:  `Setup and Run Cloud-Barista System`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n[Setup and Run Cloud-Barista]")
		fmt.Println()

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
			common.FileStr = common.GenConfigPath(common.FileStr, common.CBOperatorMode)

			var cmdStr string
			switch common.CBOperatorMode {
			case common.ModeDockerCompose:
				cmdStr = "sudo COMPOSE_PROJECT_NAME=cloud-barista docker-compose -f " + common.FileStr + " up"
				//fmt.Println(cmdStr)
				common.SysCall(cmdStr)
			case common.ModeKubernetes:
				// For Kubernetes 1.19 and above (included)
				cmdStr = "sudo kubectl create ns " + common.CBK8sNamespace + " --dry-run=client -o yaml | kubectl apply -f -"
				// For Kubernetes 1.18 and below (included)
				//cmdStr = "sudo kubectl create ns " + common.CBK8sNamespace + " --dry-run -o yaml | kubectl apply -f -"
				common.SysCall(cmdStr)
				cmdStr = "sudo helm install --namespace " + common.CBK8sNamespace + " " + common.CBHelmReleaseName + " -f " + common.FileStr + " ../helm-chart --debug"
				//fmt.Println(cmdStr)
				common.SysCall(cmdStr)
			default:

			}

		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	pf := runCmd.PersistentFlags()
	pf.StringVarP(&common.FileStr, "file", "f", common.NotDefined, "User-defined configuration file")

	/*
		switch common.CBOperatorMode {
		case common.ModeDockerCompose:
			pf.StringVarP(&common.FileStr, "file", "f", "../docker-compose-mode-files/docker-compose.yaml", "Path to Cloud-Barista Docker Compose YAML file")
		case common.ModeKubernetes:
			pf.StringVarP(&common.FileStr, "file", "f", "../helm-chart/values.yaml", "Path to Cloud-Barista Helm chart file")
		default:

		}
	*/

	//	cobra.MarkFlagRequired(pf, "file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
