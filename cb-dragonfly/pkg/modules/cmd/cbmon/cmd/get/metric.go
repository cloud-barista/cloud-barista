package get

import (
	"fmt"

	"github.com/spf13/cobra"

	pb "github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/protobuf/cbdragonfly"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/request"
)

func newGetMetricCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metric",
		Short: "Get Monitoring metric information",
		Long:  ``,
		RunE:  getMetricRun,
	}
	cmd.Flags().StringP("ns-id", "", "", "")
	cmd.Flags().StringP("mcis-id", "", "", "")
	cmd.Flags().StringP("vm-id", "", "", "")
	cmd.Flags().StringP("metric", "", "", "")
	cmd.Flags().StringP("period-type", "", "", "")
	cmd.Flags().StringP("statistics-criteria", "", "", "")
	cmd.Flags().StringP("duration", "", "", "")
	return cmd
}

func getMetricRun(cmd *cobra.Command, args []string) error {
	nsId, _ := cmd.Flags().GetString("ns-id")
	mcisId, _ := cmd.Flags().GetString("mcis-id")
	vmId, _ := cmd.Flags().GetString("vm-id")
	metricName, _ := cmd.Flags().GetString("metric")
	periodType, _ := cmd.Flags().GetString("period-type")
	statisticsCriteria, _ := cmd.Flags().GetString("statistics-criteria")
	duration, _ := cmd.Flags().GetString("duration")

	reqParams := pb.VMMonQryRequest{
		NsId:               nsId,
		McisId:             mcisId,
		VmId:               vmId,
		PeriodType:         periodType,
		StatisticsCriteria: statisticsCriteria,
		Duration:           duration,
	}

	monApi := request.GetMonitoringAPI()
	result, err := monApi.GetVMMonInfo(metricName, reqParams)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
