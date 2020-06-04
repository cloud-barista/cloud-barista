package histogram

import (
	"regexp"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics"
)

// ===== [ Constants and Variables ] =====

var (
	// Latency 패턴 정규식
	latencyPattern = core.AppName + `\.proxy\.latency\.layer\.([a-zA-Z]+)\.name\.(.*)\.complete\.(true|false)\.error\.(true|false)`
	latencyRegexp  = regexp.MustCompile(latencyPattern)

	// Router 패턴 정규식
	routerPattern = core.AppName + `\.router\.response\.(.*)\.(size|time)`
	routerRegexp  = regexp.MustCompile(routerPattern)
)

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// latencyPoints - Latency와 관련된 Metric에 대한 Histogram Points 반환
func latencyPoints(hostname string, now time.Time, histograms map[string]metrics.HistogramData, logger logging.Logger) []*client.Point {
	points := []*client.Point{}
	for key, histogram := range histograms {
		if !latencyRegexp.MatchString(key) {
			continue
		}
		if isEmpty(histogram) {
			continue
		}

		params := latencyRegexp.FindAllStringSubmatch(key, -1)[0][1:]
		tags := map[string]string{
			"host":     hostname,
			"layer":    params[0],
			"name":     params[1],
			"complete": params[2],
			"error":    params[3],
		}

		histogramPoint, err := client.NewPoint("requests", tags, newFields(histogram), now)
		if err != nil {
			logger.Error("creating histogram point:", err.Error())
			continue
		}
		points = append(points, histogramPoint)
	}
	return points
}

// routerPoints - Router와 관련된 Metric에 대한 Histogram Points 반환
func routerPoints(hostname string, now time.Time, histograms map[string]metrics.HistogramData, logger logging.Logger) []*client.Point {
	points := []*client.Point{}
	for key, histogram := range histograms {
		if !routerRegexp.MatchString(key) {
			continue
		}
		if isEmpty(histogram) {
			continue
		}

		params := routerRegexp.FindAllStringSubmatch(key, -1)[0][1:]
		tags := map[string]string{
			"host": hostname,
			"name": params[0],
		}

		histogramPoint, err := client.NewPoint("router.response-"+params[1], tags, newFields(histogram), now)
		if err != nil {
			logger.Error("creating histogram point:", err.Error())
			continue
		}
		points = append(points, histogramPoint)
	}
	return points
}

// debugPoint - Debug와 관련된 Metric에 대한 Histogram Points 반환
func debugPoint(hostname string, now time.Time, histograms map[string]metrics.HistogramData, logger logging.Logger) *client.Point {
	hd, ok := histograms[core.AppName+".service.debug.GCStats.Pause"]
	if !ok {
		return nil
	}
	tags := map[string]string{
		"host": hostname,
	}

	histogramPoint, err := client.NewPoint("service.debug.GCStats.Pause", tags, newFields(hd), now)
	if err != nil {
		logger.Error("creating histogram point:", err.Error())
		return nil
	}
	return histogramPoint
}

// runtimePoint - Runtime과 관련된 Metric에 대한 Histogram Points 반환
func runtimePoint(hostname string, now time.Time, histograms map[string]metrics.HistogramData, logger logging.Logger) *client.Point {
	hd, ok := histograms[core.AppName+".service.runtime.MemStats.PauseNs"]
	if !ok {
		return nil
	}
	tags := map[string]string{
		"host": hostname,
	}

	histogramPoint, err := client.NewPoint("service.runtime.MemStats.PauseNs", tags, newFields(hd), now)
	if err != nil {
		logger.Error("creating histogram point:", err.Error())
		return nil
	}
	return histogramPoint
}

// isEmpty - Histogram의 값이 전부 비어있는지 검사
func isEmpty(histogram metrics.HistogramData) bool {
	return histogram.Max == 0 && histogram.Min == 0 &&
		histogram.Mean == .0 && histogram.Stddev == .0 && histogram.Variance == 0 &&
		(len(histogram.Percentiles) == 0 ||
			histogram.Percentiles[0] == .0 && histogram.Percentiles[len(histogram.Percentiles)-1] == .0)
}

// newFields - 지정한 Histogram 정보를 기준으로 Percentiles 생성
func newFields(h metrics.HistogramData) map[string]interface{} {
	fields := map[string]interface{}{
		"max":      int(h.Max),
		"min":      int(h.Min),
		"mean":     int(h.Mean),
		"stddev":   int(h.Stddev),
		"variance": int(h.Variance),
	}

	if len(h.Percentiles) != 7 {
		return fields
	}

	fields["p0.1"] = h.Percentiles[0]
	fields["p0.25"] = h.Percentiles[1]
	fields["p0.5"] = h.Percentiles[2]
	fields["p0.75"] = h.Percentiles[3]
	fields["p0.9"] = h.Percentiles[4]
	fields["p0.95"] = h.Percentiles[5]
	fields["p0.99"] = h.Percentiles[6]

	return fields
}

// ===== [ Public Functions ] =====

// Points - 지정한 정보를 기준으로 InfluxDB에 적용할 Histogram points 반환
func Points(hostname string, now time.Time, histograms map[string]metrics.HistogramData, logger logging.Logger) []*client.Point {
	points := latencyPoints(hostname, now, histograms, logger)
	points = append(points, routerPoints(hostname, now, histograms, logger)...)
	if p := debugPoint(hostname, now, histograms, logger); p != nil {
		points = append(points, p)
	}
	if p := runtimePoint(hostname, now, histograms, logger); p != nil {
		points = append(points, p)
	}
	return points
}
