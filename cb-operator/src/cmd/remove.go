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
	Long:  `Stop and Remove Cloud-Barista System. Stop and Remove Cloud-Barista runtimes and related container images and meta-DB if necessary`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("\n[Remove Cloud-Barista]")
		fmt.Println()

		if common.FileStr == "" {
			fmt.Println("file is required")
		} else {
			common.FileStr = common.GenConfigPath(common.FileStr, common.CBOperatorMode)
			var cmdStr string
			switch common.CBOperatorMode {
			case common.ModeKubernetes:
				cmdStr = "sudo helm uninstall --namespace " + common.CBK8sNamespace + " " + common.CBHelmReleaseName
				common.SysCall(cmdStr)

				cmdStr = "sudo kubectl delete pvc cb-spider -n " + common.CBK8sNamespace
				common.SysCall(cmdStr)

				cmdStr = "sudo kubectl delete pvc cb-tumblebug -n " + common.CBK8sNamespace
				common.SysCall(cmdStr)

				cmdStr = "sudo kubectl delete pvc cb-ladybug -n " + common.CBK8sNamespace
				common.SysCall(cmdStr)

				cmdStr = "sudo kubectl delete pvc cb-dragonfly -n " + common.CBK8sNamespace
				common.SysCall(cmdStr)

				//fallthrough
			case common.ModeDockerCompose:
				if volFlag && imgFlag {
					cmdStr = "sudo COMPOSE_PROJECT_NAME=cloud-barista docker-compose -f " + common.FileStr + " down -v --rmi all"
				} else if volFlag {
					cmdStr = "sudo COMPOSE_PROJECT_NAME=cloud-barista docker-compose -f " + common.FileStr + " down -v"
				} else if imgFlag {
					cmdStr = "sudo COMPOSE_PROJECT_NAME=cloud-barista docker-compose -f " + common.FileStr + " down --rmi all"
				} else {
					cmdStr = "sudo COMPOSE_PROJECT_NAME=cloud-barista docker-compose -f " + common.FileStr + " down"
				}

				//fmt.Println(cmdStr)
				common.SysCall(cmdStr)

				common.SysCallDockerComposePs()
			default:

			}
		}

	},
}

var volFlag bool
var imgFlag bool

func init() {
	rootCmd.AddCommand(removeCmd)

	pf := removeCmd.PersistentFlags()
	pf.StringVarP(&common.FileStr, "file", "f", common.NotDefined, "User-defined configuration file")
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
