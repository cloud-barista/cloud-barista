package metric

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"

	"github.com/cloud-barista/cb-dragonfly/pkg/core/metric"
)

// 멀티 클라우드 인프라 서비스 개별 VM 모니터링 정보 조회
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
	result, errCode, err := metric.GetVMMonInfo(nsId, mcisId, vmId, metricName, period, aggregateType, duration)
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, result)
}
