package get

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/request"
)

func newGetConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Get configuration",
		Long:  "",
		RunE:  getConfigRun,
	}
	return cmd
}

func getConfigRun(cmd *cobra.Command, args []string) error {
	monApi := request.GetMonitoringAPI()
	result, err := monApi.GetMonitoringConfig()
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
