package cmd

import (
	"fmt"

	"github.com/cloud-barista/cb-mcks/src/grpc-api/cbadm/app"
	"github.com/cloud-barista/cb-mcks/src/utils/lang"
	"github.com/spf13/cobra"
)

type CreateClusterOptions struct {
	*app.Options
	ControlPlane struct {
		Connection string
		Count      int
		Spec       string
	}
	Worker struct {
		Connection string
		Count      int
		Spec       string
	}
}

type CreateNodeOptions struct {
	*app.Options
	clusterName string
	Worker      struct {
		Connection string
		Count      int
		Spec       string
	}
}

func (o *CreateClusterOptions) Validate() error {
	o.Namespace = lang.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
	if o.Namespace == "" {
		return fmt.Errorf("Namespace is required.")
	}
	if o.Data == "" && o.Filename == "" && o.Name == "" {
		return fmt.Errorf("One of -f Filepath or -d data is required")
	}
	return nil
}

func (o *CreateNodeOptions) Validate() error {
	o.Namespace = lang.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
	if o.Namespace == "" {
		return fmt.Errorf("Namespace is required.")
	}
	if o.clusterName == "" {
		return fmt.Errorf("ClusterName is required.")
	}
	if o.Data == "" && o.Filename == "" && o.clusterName == "" {
		return fmt.Errorf("One of -f Filepath or -d data is required")
	}
	return nil
}

func NewCreateCmd(o *app.Options) *cobra.Command {
	oCluster := &CreateClusterOptions{
		Options: o,
	}

	oNode := &CreateNodeOptions{
		Options: o,
	}

	cmds := &cobra.Command{
		Use:   "create",
		Short: "Create command",
		Long:  "This is a create command",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	cmdCluster := &cobra.Command{
		Use:   "cluster",
		Short: "Create a cluster",
		Long:  "This is a create command for cluster",
		Run: func(cmd *cobra.Command, args []string) {
			app.ValidateError(cmd, oCluster.Validate())
			app.ValidateError(cmd, func() error {
				out, err := app.GetBody(oCluster, tplCluster)
				if err != nil {
					return err
				} else {
					o.Data = `{"namespace":"` + o.Namespace + `" , "ReqInfo": ` + string(out) + `}`
				}
				SetupAndRun(cmd, o)
				return nil
			}())
		},
	}
	cmds.AddCommand(cmdCluster)
	cmdCluster.Flags().StringVar(&oCluster.ControlPlane.Connection, "control-plane-connection", "", "Connection name of control-plane nodes")
	cmdCluster.Flags().IntVar(&oCluster.ControlPlane.Count, "control-plane-count", 1, "Count of control-plane nodes")
	cmdCluster.Flags().StringVar(&oCluster.ControlPlane.Spec, "control-plane-spec", "", "Spec. of control-plane nodes")
	cmdCluster.Flags().StringVar(&oCluster.Worker.Connection, "worker-connection", "", "Connection name of wroker nodes")
	cmdCluster.Flags().IntVar(&oCluster.Worker.Count, "worker-count", 1, "Count of wroker nodes")
	cmdCluster.Flags().StringVar(&oCluster.Worker.Spec, "worker-spec", "", "Spec. of wroker nodes")

	cmdNode := &cobra.Command{
		Use:   "node (NAME | --name NAME) --cluster CLUSTER_NAME [options]",
		Short: "Create a node",
		Long:  "This is a create command for node",
		Run: func(cmd *cobra.Command, args []string) {
			app.ValidateError(cmd, oNode.Validate())
			app.ValidateError(cmd, func() error {
				out, err := app.GetBody(oNode, tplNode)
				if err != nil {
					return err
				} else {
					o.Data = `{"namespace":"` + o.Namespace + `" , "cluster":"` + oNode.clusterName + `" , "ReqInfo": ` + string(out) + `}`
				}
				SetupAndRun(cmd, o)
				return nil
			}())
		},
	}
	cmdNode.Flags().StringVar(&oNode.clusterName, "cluster", "", "Name of cluster")
	cmdNode.Flags().StringVar(&oNode.Worker.Connection, "worker-connection", "", "Connection name of wroker nodes")
	cmdNode.Flags().IntVar(&oNode.Worker.Count, "worker-count", 1, "Count of wroker nodes")
	cmdNode.Flags().StringVar(&oNode.Worker.Spec, "worker-spec", "", "Spec. of wroker nodes")
	cmds.AddCommand(cmdNode)
	/*
		cmdCredential := &cobra.Command{
			Use:   "credential",
			Short: "Create a cloud credential",
			Long:  "This is a create command for credential",
			Run: func(cmd *cobra.Command, args []string) {
				oCreate.ConvertData(cmd.Name())
				SetupAndRun(cmd, o)
			},
		}
		cmds.AddCommand(cmdCredential)
	*/
	return cmds
}

const (
	tplCluster = `{
   "name": "{{.Name}}",
   "label": "",
   "description": "",
   "controlPlane": [
      { "connection": "{{.ControlPlane.Connection}}", "count": {{.ControlPlane.Count}}, "spec": "{{.ControlPlane.Spec}}" }
   ],
   "worker": [
      { "connection": "{{.Worker.Connection}}", "count": {{.Worker.Count}}, "spec": "{{.Worker.Spec}}" }
    ],
    "config": {
        "kubernetes": {
            "networkCni": "canal",
            "podCidr": "10.244.0.0/16",
            "serviceCidr": "10.96.0.0/12",
            "serviceDnsDomain": "cluster.local"
        }
    }
}`
	tplNode = `{
	"worker": [
	   { "connection": "{{.Worker.Connection}}", "count": {{.Worker.Count}}, "spec": "{{.Worker.Spec}}" }
	 ]
}`
)
