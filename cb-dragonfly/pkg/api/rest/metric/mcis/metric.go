package mcis

import (
	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/metric"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
)

// GetVMMonInfo 멀티 클라우드 인프라 서비스 개별 VM 모니터링 정보 조회
// @Summary Get vm monitoring info
// @Description 멀티 클라우드 인프라 VM 모니터링 정보 조회
// @Tags [Monitoring] Monitoring management
// @Accept  json
// @Produce  json
// @Param ns_id path string true "네임스페이스 아이디"
// @Param mcis_id path string true "MCIS 아이디"
// @Param vm_id path string true "VM 아이디"
// @Param metric_name path string true "메트릭 정보"
// @Param periodType query string false "모니터링 단위" Enums(m, h, d)
// @Param statisticsCriteria query string false "모니터링 통계 기준" Enums(min, max, avg, last)
// @Param duration query string false "모니터링 조회 범위" Enums(5m, 5h, 5d)
// @Success 200 {object} rest.VMMonInfoType
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /ns/{ns_id}/mcis/{mcis_id}/vm/{vm_id}/metric/{metric_name}/info [get]
func GetVMMonInfo(c echo.Context) error {
	// Path 파라미터 가져오기
	nsId := c.Param("ns_id")
	mcisId := c.Param("mcis_id")
	vmId := c.Param("vm_id")
	metricName := c.Param("metric_name")
	// Query 파라미터 가져오기
	period := c.QueryParam("periodType")
	aggregateType := c.QueryParam("statisticsCriteria")
	duration := c.QueryParam("duration")
	if string(duration[len(duration)-1]) == "m" {
		durationInt, _ := strconv.Atoi(duration[:len(duration)-1])
		if durationInt < 2 {
			return echo.NewHTTPError(404, rest.SetMessage("Error! Mininum duration time is 2m"))
		}
	}

	dbInfo := types.DBMetricRequestInfo{
		NsID:                nsId,
		ServiceType:         types.MCIS,
		ServiceID:           mcisId,
		VMID:                vmId,
		MetricName:          metricName,
		MonitoringMechanism: strings.EqualFold(config.GetInstance().Monitoring.DefaultPolicy, types.PushPolicy),
		Period:              period,
		AggegateType:        aggregateType,
		Duration:            duration,
	}

	result, errCode, err := metric.GetMonInfo(dbInfo)
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, result)
}
