package get

import (
	"fmt"

	"github.com/spf13/cobra"

	pb "github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/protobuf/cbdragonfly"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/request"
)

func newGetOnDemandMetricCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ondemand-metric",
		Short: "Get OnDemand Monitoring metric information",
		Long:  ``,
		RunE:  getOnDemandMetricRun,
	}
	cmd.Flags().StringP("ns-id", "", "", "")
	cmd.Flags().StringP("mcis-id", "", "", "")
	cmd.Flags().StringP("vm-id", "", "", "")
	cmd.Flags().StringP("agent-ip", "", "", "")
	cmd.Flags().StringP("metric", "", "", "")
	return cmd
}

func getOnDemandMetricRun(cmd *cobra.Command, args []string) error {
	nsId, _ := cmd.Flags().GetString("ns-id")
	mcisId, _ := cmd.Flags().GetString("mcis-id")
	vmId, _ := cmd.Flags().GetString("vm-id")
	agentIp, _ := cmd.Flags().GetString("agent-ip")
	metricName, _ := cmd.Flags().GetString("metric")

	reqParams := pb.VMOnDemandMonQryRequest{
		NsId:    nsId,
		McisId:  mcisId,
		VmId:    vmId,
		AgentIp: agentIp,
	}

	monApi := request.GetMonitoringAPI()
	result, err := monApi.GetVMOnDemandMonInfo(metricName, reqParams)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
