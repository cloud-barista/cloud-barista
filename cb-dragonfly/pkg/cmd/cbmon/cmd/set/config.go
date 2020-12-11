package set

import (
	"fmt"

	"github.com/spf13/cobra"

	pb "github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/protobuf/cbdragonfly"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/request"
)

func newSetConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "set configuration",
		Long:  "",
		RunE:  setConfigRun,
	}
	cmd.Flags().Int32P("agent-interval", "", 10, "Set agent interval")
	cmd.Flags().Int32P("collector-interval", "", 10, "Set collector interval")
	cmd.Flags().Int32P("max-hostcnt", "", 10, "Set maximum host count")
	cmd.Flags().String("mon-policy", "", "Set mon-policy")
	return cmd
}

func setConfigRun(cmd *cobra.Command, args []string) error {
	agentInterval, _ := cmd.Flags().GetInt32("agent-interval")
	collectorInterval, _ := cmd.Flags().GetInt32("collector-interval")
	maxHostCnt, _ := cmd.Flags().GetInt32("max-hostcnt")
	monitoringPolicy, _ := cmd.Flags().GetString("mon-policy")

	reqParams := pb.MonitoringConfigInfo{
		AgentInterval:     agentInterval,
		CollectorInterval: collectorInterval,
		MaxHostCount:      maxHostCnt,
		MonitoringPolicy:  monitoringPolicy,
	}

	monApi := request.GetMonitoringAPI()
	result, err := monApi.SetMonitoringConfig(reqParams)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
