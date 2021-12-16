package cmd

import (
	"fmt"
	"strings"

	"github.com/cloud-barista/cb-operator/src/common"
	"github.com/spf13/cobra"
)

// installWeaveScopeCmd represents the install-weave-scope command
var installWeaveScopeCmd = &cobra.Command{
	Use:   "install-weave-scope",
	Short: "Install Weave Scope",
	Long:  `Install Weave Scope`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n[Install Weave Scope]")
		fmt.Println()

		common.FileStr = common.GenConfigPath(common.FileStr, common.CBOperatorMode)

		var cmdStr string
		switch common.CBOperatorMode {
		case common.ModeDockerCompose:
			fmt.Println("cb-operator Docker Compose mode does not support 'weave-scope install' subcommand.")

		case common.ModeKubernetes:
			if k8sprovider == common.NotDefined {
				fmt.Print(`--k8sprovider argument is required but not provided.
					e.g.
					--k8sprovider=gke
					--k8sprovider=eks
					--k8sprovider=aks
					--k8sprovider=mcks
					--k8sprovider=minikube
					--k8sprovider=kubeadm
					`)

				break
			}

			// If your cluster is on GKE, first you need to grant permissions for the installation.
			if strings.ToLower(k8sprovider) == "gke" {
				cmdStr = `kubectl create clusterrolebinding "cluster-admin-$(whoami)" --clusterrole=cluster-admin --user="$(gcloud config get-value core/account)"`
				common.SysCall(cmdStr)
			}

			if strings.ToLower(k8sprovider) == "gke" || strings.ToLower(k8sprovider) == "eks" || strings.ToLower(k8sprovider) == "aks" {

				// Install Weave Scope on your Kubernetes cluster.
				cmdStr = `kubectl apply -f "https://cloud.weave.works/k8s/scope.yaml?k8s-version=$(kubectl version | base64 | tr -d '\n')&k8s-service-type=LoadBalancer"`
				common.SysCall(cmdStr)

				fmt.Print(`Weave Scope installed successfully.
					To access Weave Scope UI, follow these steps:
						1. Run 'kubectl get svc -n weave'.
						2. Check EXTERNAL-IP. If EXTERNAL-IP is <pending>, then wait for seconds and run 'kubectl get svc -n weave' again.
						3. In your web browser, access to http://<EXTERNAL-IP>:80`)

			} else {
				// Install Weave Scope on your Kubernetes cluster.
				cmdStr = `kubectl apply -f "https://cloud.weave.works/k8s/scope.yaml?k8s-version=$(kubectl version | base64 | tr -d '\n')&k8s-service-type=NodePort"`
				common.SysCall(cmdStr)

				fmt.Print(`Weave Scope installed successfully.
					To access Weave Scope UI, follow these steps:
						1. Run 'kubectl get svc -n weave'.
						2. Check the NodePort port number. (30000-32768)
						3. In your web browser, access to http://<Node-IP>:<Weave-Scope-Port>`)
			}

		default:

		}

	},
}

func init() {
	weaveScopeCmd.AddCommand(installWeaveScopeCmd)

	// pf := installWeaveScopeCmd.PersistentFlags()
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
	// installWeaveScopeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installWeaveScopeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
