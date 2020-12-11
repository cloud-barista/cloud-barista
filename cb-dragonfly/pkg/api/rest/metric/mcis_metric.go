package metric

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	"github.com/cloud-barista/cb-dragonfly/pkg/core/metric"
)

// 멀티 클라우드 인트라 서비스 모니터링 메트릭 수집
func GetMCISMetric(c echo.Context) error {
	nsId := c.Param("ns_id")
	mcisId := c.Param("mcis_id")
	vmId := c.Param("vm_id")
	agentIp := c.Param("agent_ip")
	metricName := c.Param("metric_name")

	/*var mcismetric metric.MCISMetric
	if metricName == "" {
		err := errors.New("No Metric Type in API")
		return c.JSON(http.StatusInternalServerError, err)
	}
	mc.parameter.agent_ip = c.Param("agent_ip")
	// Query Agent IP 값 체크
	if mc.parameter.agent_ip == "" {
		err := errors.New("No Agent IP in API")
		return c.JSON(http.StatusInternalServerError, err)
	}
	// MCIS 모니터링 메트릭 파라미터 추출
	mc.parameter.mcis_metric = c.Param("mcis_metric_name")
	if mc.parameter.mcis_metric == "" {
		err := errors.New("No Metric Type in API")
		return c.JSON(http.StatusInternalServerError, err)
	}*/

	// MCIS 모니터링 메트릭 파라미터 기반 동작
	switch metricName {
	case "Rtt":
		rttParam := new(metric.Request)
		if err := c.Bind(rttParam); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		result, errCode, err := metric.GetMCISMonRTTInfo(nsId, mcisId, vmId, agentIp, *rttParam)
		if errCode != http.StatusOK {
			return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
		}
		return c.JSON(http.StatusOK, result)
	case "Mrtt":
		mrttParam := new(metric.Mrequest)
		if err := c.Bind(mrttParam); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		result, errCode, err := metric.GetMCISMonMRTTInfo(nsId, mcisId, vmId, agentIp, *mrttParam)
		if errCode != http.StatusOK {
			return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
		}
		return c.JSON(http.StatusOK, result)
	default:
		result, errCode, err := metric.GetMCISCommonMonInfo(nsId, mcisId, vmId, agentIp, metricName)
		if errCode != http.StatusOK {
			return echo.NewHTTPError(http.StatusInternalServerError, rest.SetMessage(err.Error()))
		}
		return c.JSON(http.StatusOK, result)
	}
	return nil
}
