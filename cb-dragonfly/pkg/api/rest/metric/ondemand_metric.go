package metric

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/core/metric"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
)

// GetVMOnDemandMetric 멀티클라우드 인프라 VM 온디멘드 모니터링
// @Summary Get vm on-demand monitoring metric info
// @Description 멀티 클라우드 인프라 VM 온디맨드 모니터링 정보 조회
// @Tags [Monitoring] Monitoring management
// @Accept  json
// @Produce  json
// @Param ns path string true "네임스페이스 아이디"
// @Param mcis_id path string true "MCIS 아이디"
// @Param vm_id path string true "VM 아이디"
// @Param agent_ip path string true "에이전트 아이피"
// @Param metric_name path string true "메트릭 정보"
// @Success 200 {object} rest.VMOnDemandMetricType
// @Failure 404 {object} rest.SimpleMsg
// @Failure 500 {object} rest.SimpleMsg
// @Router /ns/{ns}/mcis/{mcis_id}/vm/{vm_id}/agent_ip/{agent_ip}/metric/{metric_name}/ondemand-monitoring-info [get]
func GetVMOnDemandMetric(c echo.Context) error {
	// 온디멘드 모니터링 Path 파라미터 가져오기
	//nsId := c.Param("ns_id")
	//mcisId := c.Param("mcis_id")
	//vmId := c.Param("vm_id")
	publicIP := c.Param("agent_ip")
	metricName := c.Param("metric_name")

	// 파라미터 값 체크
	if publicIP == "" || metricName == "" {
		return c.JSON(http.StatusInternalServerError, errors.New("parameter is missing"))
	}

	result, errCode, err := metric.GetVMOnDemandMonInfo(metricName, publicIP)
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, result)
}

func GetMCISOnDemandPacket(c echo.Context) error {
	// 온디멘드 모니터링 Path 파라미터 가져오기
	nsId := c.Param("ns_id")
	mcisId := c.Param("mcis_id")
	vmId := c.Param("vm_id")
	watchTime := c.Param("watch_time")

	// 파라미터 값 체크
	if nsId == "" || mcisId == "" || vmId == "" {
		return c.JSON(http.StatusInternalServerError, errors.New("parameter is missing"))
	}

	result, errCode, err := metric.GetMCISOnDemandPacketInfo(nsId, mcisId, vmId, watchTime)
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, result)
}

func GetMCISOnDemandProcess(c echo.Context) error {
	// 온디멘드 모니터링 Path 파라미터 가져오기
	publicIP := c.Param("agent_ip")

	// 파라미터 값 체크
	if publicIP == "" {
		return c.JSON(http.StatusInternalServerError, errors.New("parameter is missing"))
	}

	result, errCode, err := metric.GetMCISOnDemandProcessInfo(publicIP)
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, result)
}

func GetMCISSpec(c echo.Context) error {
	// 온디멘드 모니터링 Path 파라미터 가져오기
	nsId := c.Param("ns")
	mcisId := c.Param("mcis_id")
	auth := c.Request().Header.Get("Authorization")
	// 파라미터 값 체크
	if nsId == "" || mcisId == "" {
		return c.JSON(http.StatusInternalServerError, errors.New("parameter is missing"))
	}
	if auth == "" {
		return c.JSON(http.StatusInternalServerError, errors.New("basic auth is missing"))
	}

	result, errCode, err2 := metric.GetMCISSpecInfo(nsId, mcisId, auth)
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err2.Error()))
	}
	return c.JSON(http.StatusOK, result)
}
