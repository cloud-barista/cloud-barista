package cmd

import (
	"fmt"

	"github.com/cloud-barista/cb-operator/src/common"
	"github.com/spf13/cobra"
)

// weaveScopeCmd represents the weave-scope command
var weaveScopeCmd = &cobra.Command{
	Use:   "weave-scope",
	Short: "Subcommand for managing Weave Scope",
	Long:  `Subcommand for managing Weave Scope`,
	Run: func(cmd *cobra.Command, args []string) {

		switch common.CBOperatorMode {
		case common.ModeDockerCompose:
			fmt.Println("cb-operator Docker Compose mode does not support 'uninstall-weave-scope' subcommand.")

		case common.ModeKubernetes:

			fmt.Println("")
			fmt.Println("'./operator weave-scope' subcommand provides these subsubcommands:")
			fmt.Println("")
			fmt.Println("'./operator weave-scope install': Install and expose Weave Scope on your K8s cluster.")
			fmt.Println("'./operator weave-scope uninstall': Uninstall Weave Scope on your K8s cluster.")
			fmt.Println("")

		default:

		}
	},
}

func init() {
	rootCmd.AddCommand(weaveScopeCmd)

	pf := weaveScopeCmd.PersistentFlags()
	// pf.StringVarP(&common.FileStr, "file", "f", common.NotDefined, "User-defined configuration file")
	pf.StringVarP(&k8sprovider, "k8sprovider", "", common.NotDefined, "Kind of Managed K8s services")

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
	// weaveScopeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// weaveScopeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
