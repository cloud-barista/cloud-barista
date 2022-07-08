package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cloud-barista/cb-mcks/src/grpc-api/cbadm/app"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"
)

func NewGetCmd(o *app.Options) *cobra.Command {

	fnValidate := func() error {
		o.Namespace = lang.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
		if o.Namespace == "" {
			return fmt.Errorf("Namespace is required.")
		}
		return nil
	}

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get command",
		Long:  "This is a get command",
	}

	getCmd.AddCommand(&cobra.Command{
		Use:   "cluster (NAME | --name NAME) [options]",
		Short: "Get cluster or cluster list",
		Long:  "This is a get command for cluster",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(cmd *cobra.Command, args []string) {
			app.ValidateError(cmd, fnValidate())

			SetupAndRun(cmd, o)
		},
	})
	cmdNode := &cobra.Command{
		Use:   "node (NAME | --name NAME) [options]",
		Short: "Get node or node list",
		Long:  "This is a get command for node",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(cmd *cobra.Command, args []string) {
			app.ValidateError(cmd, fnValidate())
			app.ValidateError(cmd, func() error {
				if clusterName == "" {
					return fmt.Errorf("cluster name is required")
				}
				return nil
			}())
			SetupAndRun(cmd, o)
		},
	}
	cmdNode.Flags().StringVar(&clusterName, "cluster", "", "Name of cluster")
	getCmd.AddCommand(cmdNode)
	/*
		getCmd.AddCommand(&cobra.Command{
			Use:   "credential (NAME | --name NAME) [options]",
			Short: "Get credential",
			Long:  "This is get command for credential",
			Run: func(cmd *cobra.Command, args []string) {
				SetupAndRun(cmd, o)
			},
		})
	*/
	return getCmd
}
