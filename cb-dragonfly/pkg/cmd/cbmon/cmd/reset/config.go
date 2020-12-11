package reset

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/request"
)

func newResetConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Reset configuration",
		Long:  "",
		RunE:  resetConfigRun,
	}
	return cmd
}

func resetConfigRun(cmd *cobra.Command, args []string) error {
	monApi := request.GetMonitoringAPI()
	result, err := monApi.ResetMonitoringConfig()
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
