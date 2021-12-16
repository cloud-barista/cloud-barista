package metric

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/metric"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
)

// GetMCISMetric 멀티 클라우드 인트라 서비스 모니터링 메트릭 수집
// @Summary Get MCIS on-demand monitoring metric info
// @Description 멀티 클라우드 인프라 VM 온디맨드 모니터링 정보 조회
// @Tags [Monitoring] Monitoring management
// @Accept  json
// @Produce  json
// @Param ns_id path string true "네임스페이스 아이디"
// @Param mcis_id path string true "MCIS 아이디"
// @Param vm_id path string true "VM 아이디"
// @Param agent_ip path string true "에이전트 아이피"
// @Param metric_name path string true "메트릭 정보"
// @Success 200 {object} rest.JSONResult{[DEFAULT]=CBMCISMetric,[Mrtt]=MCBMCISMetric} "Different return structures by the given param"
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /ns/{ns_id}/mcis/{mcis_id}/vm/{vm_id}/agent_ip/{agent_ip}/mcis_metric/{metric_name}/mcis-monitoring-info [get]
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
