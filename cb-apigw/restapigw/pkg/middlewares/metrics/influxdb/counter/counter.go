package counter

import (
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/influxdata/influxdb/client/v2"
)

// ===== [ Constants and Variables ] =====

var (
	// 최종 Request 수 관리
	lastRequestCount = map[string]int{}
	// 최종 Response 수 관리
	lastResponseCount = map[string]int{}

	// Lock 운영
	mu = new(sync.Mutex)

	// Request 패턴 정규 표현식
	requestCounterPattern = core.AppName + `\.proxy\.requests\.layer\.([a-zA-Z]+)\.name\.(.*)\.complete\.(true|false)\.error\.(true|false)`
	requestCounterRegexp  = regexp.MustCompile(requestCounterPattern)
	// Response 패턴 정규 표현식
	responseCounterPattern = core.AppName + `\.router\.response\.(.*)\.status\.([\d]{3})\.count`
	responseCounterRegexp  = regexp.MustCompile(responseCounterPattern)
)

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// requestPoints - Request와 관련된 Metric에 대한 Points 반환
func requestPoints(hostname string, now time.Time, counters map[string]int64, logger logging.Logger) []*client.Point {
	points := []*client.Point{}
	mu.Lock()
	for key, count := range counters {
		if !requestCounterRegexp.MatchString(key) {
			continue
		}
		params := requestCounterRegexp.FindAllStringSubmatch(key, -1)[0][1:]
		tags := map[string]string{
			"host":     hostname,
			"layer":    params[0],
			"name":     params[1],
			"complete": params[2],
			"error":    params[3],
		}
		last, ok := lastRequestCount[strings.Join(params, ".")]
		if !ok {
			last = 0
		}
		fields := map[string]interface{}{
			"total": int(count),
			"count": int(count) - last,
		}
		lastRequestCount[strings.Join(params, ".")] = int(count)

		countersPoint, err := client.NewPoint("requests", tags, fields, now)
		if nil != err {
			logger.Error("[METRICS] InfluxDB > Creating request counters point:", err.Error())
			continue
		}

		points = append(points, countersPoint)
	}

	mu.Unlock()
	return points
}

// responsePoints - Response와 관련된 Metric에 대한 Points 반환
func responsePoints(hostname string, now time.Time, counters map[string]int64, logger logging.Logger) []*client.Point {
	points := []*client.Point{}
	mu.Lock()
	for key, count := range counters {
		if !responseCounterRegexp.MatchString(key) {
			continue
		}
		params := responseCounterRegexp.FindAllStringSubmatch(key, -1)[0][1:]
		tags := map[string]string{
			"host":   hostname,
			"name":   params[0],
			"status": params[1],
		}
		last, ok := lastResponseCount[strings.Join(params, ".")]
		if !ok {
			last = 0
		}
		fields := map[string]interface{}{
			"total": int(count),
			"count": int(count) - last,
		}
		lastResponseCount[strings.Join(params, ".")] = int(count)

		countersPoint, err := client.NewPoint("responses", tags, fields, now)
		if nil != err {
			logger.Error("[METRICS] InfluxDB > Creating response counters point:", err.Error())
			continue
		}
		points = append(points, countersPoint)
	}
	mu.Unlock()
	return points
}

// connectionPoints - Connection/Disconnection에 관련된 Metric에 대한 Points 반환
func connectionPoints(hostname string, now time.Time, counters map[string]int64, logger logging.Logger) []*client.Point {
	points := make([]*client.Point, 2)

	prefix := core.AppName + ".router."

	in := map[string]interface{}{
		"current": int(counters[prefix+"connected"]),
		"total":   int(counters[prefix+"connected-total"]),
	}
	incoming, err := client.NewPoint("router", map[string]string{"host": hostname, "direction": "in"}, in, now)
	if nil != err {
		logger.Error("[METRICS] InfluxDB > Creating incoming connection counters point:", err.Error())
		return points
	}
	points[0] = incoming

	out := map[string]interface{}{
		"current": int(counters[prefix+"disconnected"]),
		"total":   int(counters[prefix+"disconnected-total"]),
	}
	outgoing, err := client.NewPoint("router", map[string]string{"host": hostname, "direction": "out"}, out, now)
	if nil != err {
		logger.Error("[METRICS] InfluxDB > Creating outgoing connection counters point:", err.Error())
		return points
	}
	points[1] = outgoing

	return points
}

// ===== [ Public Functions ] =====

// Points - 지정한 정보를 기준으로 InfluxDB에 적용할 Counter points 반환
func Points(hostname string, now time.Time, counters map[string]int64, logger logging.Logger) []*client.Point {
	points := requestPoints(hostname, now, counters, logger)
	points = append(points, responsePoints(hostname, now, counters, logger)...)
	points = append(points, connectionPoints(hostname, now, counters, logger)...)
	return points
}
