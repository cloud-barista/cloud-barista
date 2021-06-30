package cmd

import (
	"fmt"

	"github.com/cloud-barista/cb-operator/src/common"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information of Cloud-Barista System",
	Long:  `Get information of Cloud-Barista System. Information about containers and container images`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n[Get info for Cloud-Barista runtimes]")
		fmt.Println()

		if common.FileStr == "" {
			fmt.Println("file is required")
		} else {
			common.FileStr = common.GenConfigPath(common.FileStr, common.CBOperatorMode)
			var cmdStr string
			switch common.CBOperatorMode {
			case common.ModeDockerCompose:
				common.SysCallDockerComposePs()

				fmt.Println("")
				fmt.Println("[v]Status of Cloud-Barista runtime images")
				cmdStr = "sudo COMPOSE_PROJECT_NAME=cloud-barista docker-compose -f " + common.FileStr + " images"
				//fmt.Println(cmdStr)
				common.SysCall(cmdStr)
			case common.ModeKubernetes:
				fmt.Println("[v]Status of Cloud-Barista Helm release")
				cmdStr = "sudo helm status --namespace " + common.CBK8sNamespace + " " + common.CBHelmReleaseName
				common.SysCall(cmdStr)
				fmt.Println()
				fmt.Println("[v]Status of Cloud-Barista pods")
				cmdStr = "sudo kubectl get pods -n " + common.CBK8sNamespace
				common.SysCall(cmdStr)
				fmt.Println()
				fmt.Println("[v]Status of Cloud-Barista container images")
				cmdStr = `sudo kubectl get pods -n ` + common.CBK8sNamespace + ` -o jsonpath="{..image}" |\
				tr -s '[[:space:]]' '\n' |\
				sort |\
				uniq`
				common.SysCall(cmdStr)
			default:

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	pf := infoCmd.PersistentFlags()
	pf.StringVarP(&common.FileStr, "file", "f", common.NotDefined, "User-defined configuration file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
