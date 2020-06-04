package gauge

import (
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// Points - 지정한 정보를 기준으로 InfluxDB에 적용할 Gauge points 반환
func Points(hostname string, now time.Time, counters map[string]int64, logger logging.Logger) []*client.Point {
	points := make([]*client.Point, 4)

	prefix := core.AppName + ".router."
	prefixS := core.AppName + ".service."

	in := map[string]interface{}{
		"gauge": int(counters[prefix+"connected-gauge"]),
	}
	incoming, err := client.NewPoint("router", map[string]string{"host": hostname, "direction": "in"}, in, now)
	if err != nil {
		logger.Error("creating incoming connection counters point:", err.Error())
		return points
	}
	points[0] = incoming

	out := map[string]interface{}{
		"guage": int(counters[prefix+"disconnected-gauge"]),
	}
	outgoing, err := client.NewPoint("router", map[string]string{"host": hostname, "direction": "out"}, out, now)
	if err != nil {
		logger.Error("creating outgoing connection counters point:", err.Error())
		return points
	}
	points[1] = outgoing

	debug := map[string]interface{}{}
	runtime := map[string]interface{}{}

	for key, v := range counters {
		if strings.HasPrefix(key, prefix+"connected-gauge") || strings.HasPrefix(key, prefix+"disconnected-gauge") {
			continue
		}
		if strings.HasPrefix(key, prefixS+"debug.") {
			debug[key[len(prefixS+"debug."):]] = int(v)
			continue
		}
		if strings.HasPrefix(key, prefixS+"runtime.") {
			runtime[key[len(prefixS+"runtime."):]] = int(v)
			continue
		}
		logger.Warn("unknown gauge key:", key)
	}

	debugPoint, err := client.NewPoint("debug", map[string]string{"host": hostname}, debug, now)
	if err != nil {
		logger.Error("creating debug counters point:", err.Error())
		return points
	}
	points[2] = debugPoint

	runtimePoint, err := client.NewPoint("runtime", map[string]string{"host": hostname}, runtime, now)
	if err != nil {
		logger.Error("creating runtime counters point:", err.Error())
		return points
	}
	points[3] = runtimePoint

	return points
}
