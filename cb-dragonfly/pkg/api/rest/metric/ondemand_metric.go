package metric

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/rest"
	"github.com/cloud-barista/cb-dragonfly/pkg/core/metric"
)

// 멀티클라우드 인프라 VM 온디멘드 모니터링
func GetVMOnDemandMetric(c echo.Context) error {
	// 온디멘드 모니터링 Path 파라미터 가져오기
	nsId := c.Param("ns_id")
	mcisId := c.Param("mcis_id")
	vmId := c.Param("vm_id")
	publicIP := c.Param("agent_ip")
	metricName := c.Param("metric_name")

	// 파라미터 값 체크
	if publicIP == "" || metricName == "" {
		return c.JSON(http.StatusInternalServerError, errors.New("parameter is missing"))
	}

	result, errCode, err := metric.GetVMOnDemandMonInfo(nsId, mcisId, vmId, metricName, publicIP)
	if errCode != http.StatusOK {
		return echo.NewHTTPError(errCode, rest.SetMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, result)
}

/*
func MappingMonMetric(metricKey string, metricVal collector.TelegrafMetric) (map[string]interface{}, error) {
	metricMap := map[string]interface{}{}
	metricMap["name"] = metricVal.Name
	tagMap := map[string]interface{}{
		"nsId":   metricVal.Tags["nsId"],
		"mcisId": metricVal.Tags["mcisId"],
		"vmId":   metricVal.Tags["vmId"],
	}
	metricMap["tags"] = tagMap

	metricCols, err := influxdb.MappingOnDemandMonMetric(metricKey, metricVal.Fields)
	if err != nil {
		return nil, err
	}
	metricMap["values"] = metricCols
	metricMap["time"] = time.Now().UTC() // TODO: parsing timestamp to utc time
	return metricMap, nil
}
*/
