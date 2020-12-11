package get

import (
	"fmt"

	"github.com/spf13/cobra"

	pb "github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/protobuf/cbdragonfly"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/request"
)

func newGetMCISMetricCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcis-metric",
		Short: "Get MCIS Monitoring metric information",
		Long:  ``,
		RunE:  getMCISMetricRun,
	}
	cmd.Flags().StringP("ns-id", "", "", "")
	cmd.Flags().StringP("mcis-id", "", "", "")
	cmd.Flags().StringP("vm-id", "", "", "")
	cmd.Flags().StringP("agent-ip", "", "", "")
	cmd.Flags().StringP("metric", "", "", "")
	return cmd
}

func getMCISMetricRun(cmd *cobra.Command, args []string) error {
	nsId, _ := cmd.Flags().GetString("ns-id")
	mcisId, _ := cmd.Flags().GetString("mcis-id")
	vmId, _ := cmd.Flags().GetString("vm-id")
	agentIp, _ := cmd.Flags().GetString("agent-ip")
	metricName, _ := cmd.Flags().GetString("metric")

	reqParams := pb.VMMCISMonQryRequest{
		NsId:       nsId,
		McisId:     mcisId,
		VmId:       vmId,
		AgentIp:    agentIp,
		MetricName: metricName,
	}

	monApi := request.GetMonitoringAPI()
	result, err := monApi.GetMCISMonInfo(reqParams)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
