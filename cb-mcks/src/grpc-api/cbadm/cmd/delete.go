package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cloud-barista/cb-mcks/src/grpc-api/cbadm/app"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"
)

// returns a cobra command
func NewDeleteCmd(o *app.Options) *cobra.Command {

	fnValidate := func() error {
		o.Namespace = lang.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
		if o.Namespace == "" {
			return fmt.Errorf("Namespace is required.")
		}
		if o.Name == "" {
			return fmt.Errorf("Name is required.")
		}
		return nil
	}

	// root
	cmds := &cobra.Command{
		Use:   "delete",
		Short: "Delete command",
		Long:  "This is a delete command",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	// cluster
	cmds.AddCommand(&cobra.Command{
		Use:   "cluster (NAME | --name NAME) [options]",
		Short: "Delete a cluster",
		Long:  "This is a delete command for cluster",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(cmd *cobra.Command, args []string) {
			app.ValidateError(cmd, fnValidate())

			SetupAndRun(cmd, o)
		},
	})

	// node
	cmdNode := &cobra.Command{
		Use:   "node (NAME | --name NAME) --cluster CLUSTER_NAME [options]",
		Short: "Delete a node",
		Long:  "This is a delete command for node",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(cmd *cobra.Command, args []string) {
			app.ValidateError(cmd, fnValidate())
			app.ValidateError(cmd, func() error {
				if clusterName == "" {
					return fmt.Errorf("ClusterName is required")
				}
				return nil
			}())
			SetupAndRun(cmd, o)

		},
	}
	cmdNode.Flags().StringVar(&clusterName, "cluster", "", "Name of cluster")
	cmds.AddCommand(cmdNode)

	// credential
	/*
		cmds.AddCommand(&cobra.Command{
			Use:   "credential (NAME | --name NAME) [options]",
			Short: "Delete a cloud credential",
			Long:  "This is a delete command for credential",
			Args:  app.BindCommandArgs(&o.Name),
			Run: func(cmd *cobra.Command, args []string) {
				app.ValidateError(cmd, func() error {
					if o.Name == "" {
						return fmt.Errorf("Name is required.")
					}
					return nil
				}())
				SetupAndRun(cmd, o)
			},
		})*/

	return cmds
}
