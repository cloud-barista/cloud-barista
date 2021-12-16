package get

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",
		//Short: "Get a controller or resource to the project",
		//Long:  "",
	}
	cmd.AddCommand(newGetConfigCmd())
	cmd.AddCommand(newGetMetricCmd())
	cmd.AddCommand(newGetOnDemandMetricCmd())
	cmd.AddCommand(newGetMCISMetricCmd())
	return cmd
}
