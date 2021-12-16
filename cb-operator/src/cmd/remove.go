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
				cmdStr = fmt.Sprintf("helm uninstall --namespace %s %s", common.CBK8sNamespace, common.CBHelmReleaseName)
				common.SysCall(cmdStr)

				cmdStr = fmt.Sprintf("kubectl delete pvc cb-spider -n %s", common.CBK8sNamespace)
				common.SysCall(cmdStr)

				cmdStr = fmt.Sprintf("kubectl delete pvc cb-tumblebug -n %s", common.CBK8sNamespace)
				common.SysCall(cmdStr)

				cmdStr = fmt.Sprintf("kubectl delete pvc cb-ladybug -n %s", common.CBK8sNamespace)
				common.SysCall(cmdStr)

				cmdStr = fmt.Sprintf("kubectl delete pvc cb-dragonfly -n %s", common.CBK8sNamespace)
				common.SysCall(cmdStr)

				cmdStr = fmt.Sprintf("kubectl delete pvc data-cloud-barista-etcd-0 -n %s", common.CBK8sNamespace)
				common.SysCall(cmdStr)

				//fallthrough
			case common.ModeDockerCompose:
				if volFlag && imgFlag {
					cmdStr = fmt.Sprintf("COMPOSE_PROJECT_NAME=%s docker-compose -f %s down -v --rmi all", common.CBComposeProjectName, common.FileStr)
				} else if volFlag {
					cmdStr = fmt.Sprintf("COMPOSE_PROJECT_NAME=%s docker-compose -f %s down -v", common.CBComposeProjectName, common.FileStr)
				} else if imgFlag {
					cmdStr = fmt.Sprintf("COMPOSE_PROJECT_NAME=%s docker-compose -f %s down --rmi all", common.CBComposeProjectName, common.FileStr)
				} else {
					cmdStr = fmt.Sprintf("COMPOSE_PROJECT_NAME=%s docker-compose -f %s down", common.CBComposeProjectName, common.FileStr)
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
