package v1

import (
	"errors"
	"fmt"
	"time"

	influxBuilder "github.com/Scalingo/go-utils/influx"
)

func BuildQuery(isPush bool, vmId string, metric string, period string, aggregateType string, duration string) (string, error) {

	// 통계 기준 설정
	if aggregateType == "avg" {
		aggregateType = "mean"
	}

	// 시간 범위 설정
	timeDuration := fmt.Sprintf("(now()+1m) - %s", duration)

	// 시간 단위 설정
	var timeCriteria time.Duration

	// InfluXDB 쿼리 생성
	var query influxBuilder.Query

	switch metric {

	case "cpu":
		query = influxBuilder.NewQuery().On(metric).
			Field("cpu_utilization", aggregateType).
			Field("cpu_system", aggregateType).
			Field("cpu_idle", aggregateType).
			Field("cpu_iowait", aggregateType).
			Field("cpu_hintr", aggregateType).
			Field("cpu_sintr", aggregateType).
			Field("cpu_user", aggregateType).
			Field("cpu_nice", aggregateType).
			Field("cpu_steal", aggregateType).
			Field("cpu_guest", aggregateType).
			Field("cpu_guest_nice", aggregateType)

	case "cpufreq":
		query = influxBuilder.NewQuery().On(metric).
			Field("cpu_speed", aggregateType)

	case "mem":
		query = influxBuilder.NewQuery().On(metric).
			Field("mem_utilization", aggregateType).
			Field("mem_total", aggregateType).
			Field("mem_used", aggregateType).
			Field("mem_free", aggregateType).
			Field("mem_shared", aggregateType).
			Field("mem_buffers", aggregateType).
			Field("mem_cached", aggregateType)

	case "disk":
		query = influxBuilder.NewQuery().On(metric).
			Field("disk_utilization", aggregateType).
			Field("disk_total", aggregateType).
			Field("disk_used", aggregateType).
			Field("disk_free", aggregateType)

	case "diskio":
		fieldArr := []string{"kb_read", "kb_written", "ops_read", "ops_write", "read_time", "write_time"}
		query := getPerSecMetric(isPush, vmId, metric, period, fieldArr, duration)
		return query, nil

	case "net":
		fieldArr := []string{"bytes_in", "bytes_out", "pkts_in", "pkts_out", "err_in", "err_out", "drop_in", "drop_out"}
		query := getPerSecMetric(isPush, vmId, metric, period, fieldArr, duration)
		return query, nil

	default:
		return "", errors.New("not found metric")
	}

	if isPush {
		switch period {
		case "m":
			timeCriteria = time.Minute
		case "h":
			timeCriteria = time.Hour
		case "d":
			timeCriteria = time.Hour * 24
		}
		query = query.Where("time", influxBuilder.MoreThan, timeDuration).
			And("\"vmId\"", influxBuilder.Equal, "'"+vmId+"'").
			GroupByTime(timeCriteria).
			GroupByTag("\"vmId\"").
			Fill("0").
			OrderByTime("ASC")
	} else {
		query = query.Where("time", influxBuilder.MoreThan, timeDuration).
			And("\"vmId\"", influxBuilder.Equal, "'"+vmId+"'").
			GroupByTag("\"vmId\"").
			GroupByTag("\"nsId\"").
			GroupByTag("\"mcisId\"").
			//GroupByTime(timeCriteria).
			Fill("0").
			OrderByTime("ASC")
	}
	queryString := query.Build()

	return queryString, nil
}

func getPerSecMetric(isPUSH bool, vmId, metric, period string, fieldArr []string, duration string) string {
	var query string

	var timeCriteria string
	switch period {
	case "m":
		timeCriteria = "1m"
	case "h":
		timeCriteria = "1h"
	case "d":
		timeCriteria = "24h"
	}

	// 메트릭 필드 조회 쿼리 생성
	fieldQueryForm := " non_negative_derivative(first(%s), 1s) as \"%s\""
	for idx, field := range fieldArr {
		if idx == 0 {
			query = "SELECT"
		}
		query += fmt.Sprintf(fieldQueryForm, field, field)
		if idx != len(fieldArr)-1 {
			query += ","
		}
	}
	var whereQueryForm string

	// 메트릭 조회 조건 쿼리 생성
	if isPUSH {
		whereQueryForm = " FROM \"%s\" WHERE time > (now()+1m) - %s AND \"vmId\"='%s' GROUP BY time(%s) fill(0)"
		query += fmt.Sprintf(whereQueryForm, metric, duration, vmId, timeCriteria)
	} else {
		whereQueryForm = " FROM \"%s\" WHERE time > (now()+1m) - %s AND \"vmId\"='%s' GROUP BY time(%s), \"vmId\", \"nsId\", \"mcisId\" fill(0)"
		query += fmt.Sprintf(whereQueryForm, metric, duration, vmId, timeCriteria)
	}

	return query
}
