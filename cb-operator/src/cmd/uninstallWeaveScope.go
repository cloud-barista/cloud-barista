package cmd

import (
	"fmt"
	"strings"

	"github.com/cloud-barista/cb-operator/src/common"
	"github.com/spf13/cobra"
)

// uninstallWeaveScopeCmd represents the uninstall-weave-scope command
var uninstallWeaveScopeCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Weave Scope",
	Long:  `Uninstall Weave Scope`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n[Uninstall Weave Scope]")
		fmt.Println()

		common.FileStr = common.GenConfigPath(common.FileStr, common.CBOperatorMode)

		var cmdStr string
		switch common.CBOperatorMode {
		case common.ModeDockerCompose:
			fmt.Println("cb-operator Docker Compose mode does not support 'weave-scope uninstall' subcommand.")

		case common.ModeKubernetes:
			// If your cluster is on GKE, first you need to grant permissions for the uninstallation.
			if strings.ToLower(k8sprovider) == "gke" {
				cmdStr = `kubectl delete clusterrolebinding "cluster-admin-$(whoami)"`
				common.SysCall(cmdStr)
			}

			if strings.ToLower(k8sprovider) == "gke" || strings.ToLower(k8sprovider) == "eks" || strings.ToLower(k8sprovider) == "aks" {

				// Uninstall Weave Scope on your Kubernetes cluster.
				cmdStr = `kubectl delete -f "https://cloud.weave.works/k8s/scope.yaml?k8s-version=$(kubectl version | base64 | tr -d '\n')&k8s-service-type=LoadBalancer"`
				common.SysCall(cmdStr)

				fmt.Print(`Weave Scope uninstalled successfully.`)

			} else {
				// Uninstall Weave Scope on your Kubernetes cluster.
				cmdStr = `kubectl delete -f "https://cloud.weave.works/k8s/scope.yaml?k8s-version=$(kubectl version | base64 | tr -d '\n')&k8s-service-type=NodePort"`
				common.SysCall(cmdStr)

				fmt.Print(`Weave Scope uninstalled successfully.`)
			}

		default:

		}

	},
}

func init() {
	weaveScopeCmd.AddCommand(uninstallWeaveScopeCmd)

	// pf := uninstallWeaveScopeCmd.PersistentFlags()
	// // pf.StringVarP(&common.FileStr, "file", "f", common.NotDefined, "User-defined configuration file")
	// pf.StringVarP(&k8sprovider, "k8sprovider", "", common.NotDefined, "Kind of Managed K8s services")

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
	// uninstallWeaveScopeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallWeaveScopeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
