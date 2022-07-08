package v1

import (
	"errors"
	"fmt"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"strings"
	"time"

	influxBuilder "github.com/Scalingo/go-utils/influx"
)

func BuildQuery(info types.DBMetricRequestInfo) (string, error) {
	// 평균 InfluxQL 기준으로 변경
	if info.AggegateType == "avg" {
		info.AggegateType = "mean"
	}
	mcisType := util.CheckMCISType(info.ServiceType)
	mck8sType := util.CheckMCK8SType(info.ServiceType)
	// 시간 범위 설정
	timeDuration := fmt.Sprintf("(now()+1m) - %s", info.Duration)

	// 시간 단위 설정
	var timeCriteria time.Duration

	// InfluXDB 쿼리 생성
	var query influxBuilder.Query
	var diskQuery string
	// MCIS 모니터링
	if mcisType {
		switch info.MetricName {

		case "cpu":
			query = influxBuilder.NewQuery().On(info.MetricName).
				Field("cpu_utilization", info.AggegateType).
				Field("cpu_system", info.AggegateType).
				Field("cpu_idle", info.AggegateType).
				Field("cpu_iowait", info.AggegateType).
				Field("cpu_hintr", info.AggegateType).
				Field("cpu_sintr", info.AggegateType).
				Field("cpu_user", info.AggegateType).
				Field("cpu_nice", info.AggegateType).
				Field("cpu_steal", info.AggegateType).
				Field("cpu_guest", info.AggegateType).
				Field("cpu_guest_nice", info.AggegateType)

		case "cpufreq":
			query = influxBuilder.NewQuery().On(info.MetricName).
				Field("cpu_speed", info.AggegateType)

		case "mem":
			query = influxBuilder.NewQuery().On(info.MetricName).
				Field("mem_utilization", info.AggegateType).
				Field("mem_total", info.AggegateType).
				Field("mem_used", info.AggegateType).
				Field("mem_free", info.AggegateType).
				Field("mem_shared", info.AggegateType).
				Field("mem_buffers", info.AggegateType).
				Field("mem_cached", info.AggegateType)

		case "disk":
			query = influxBuilder.NewQuery().On(info.MetricName).
				Field("disk_utilization", info.AggegateType).
				Field("disk_total", info.AggegateType).
				Field("disk_used", info.AggegateType).
				Field("disk_free", info.AggegateType)

		case "diskio":
			diskQuery = getPerSecMetric(info, "kb_read", "kb_written", "ops_read", "ops_write", "read_time", "write_time")
			return diskQuery, nil

		case "net":
			diskQuery = getPerSecMetric(info, "bytes_in", "bytes_out", "pkts_in", "pkts_out", "err_in", "err_out", "drop_in", "drop_out")
			return diskQuery, nil

		default:
			return "", errors.New("not found metric")
		}
	}

	// MCK8S 모니터링
	if mck8sType {
		switch info.MetricName {
		case "kubernetes_node":
			query = influxBuilder.NewQuery().On(info.MetricName).
				Field("cpu_usage_core_nanoseconds", info.AggegateType).
				Field("memory_usage_bytes", info.AggegateType).
				Field("memory_available_bytes", info.AggegateType).
				Field("network_rx_bytes", info.AggegateType).
				Field("network_rx_errors", info.AggegateType).
				Field("network_tx_bytes", info.AggegateType).
				Field("network_tx_errors", info.AggegateType).
				Field("fs_capacity_bytes", info.AggegateType).
				Field("fs_used_bytes", info.AggegateType)

		case "kubernetes_pod_container":
			query = influxBuilder.NewQuery().On(info.MetricName).
				Field("cpu_usage_nanocores", info.AggegateType).
				Field("memory_usage_bytes", info.AggegateType).
				Field("rootfs_capacity_bytes", info.AggegateType).
				Field("rootfs_used_bytes", info.AggegateType)

		case "kubernetes_pod_network":
			networkQuery := getPerSecMetric(info, "rx_bytes", "rx_errors", "tx_bytes", "tx_errors")
			return networkQuery, nil

		default:
			return "", errors.New("not found metric")
		}
	}

	if info.MonitoringMechanism {
		switch info.Period {
		case "m":
			timeCriteria = time.Minute
		case "h":
			timeCriteria = time.Hour
		case "d":
			timeCriteria = time.Hour * 24
		}
		if mcisType {
			query = query.Where("time", influxBuilder.MoreThan, timeDuration).
				And("\"vmId\"", influxBuilder.Equal, "'"+info.VMID+"'").
				GroupByTime(timeCriteria).
				GroupByTag("\"vmId\"").
				Fill("0").
				OrderByTime("ASC")
		}
		if mck8sType {
			switch info.MetricName {
			case "kubernetes_node":
				if strings.EqualFold(info.MCK8SReqInfo.GroupBy, types.Cluster) {
					query = query.Where("time", influxBuilder.MoreThan, timeDuration).
						And("\"nsId\"", influxBuilder.Equal, "'"+info.NsID+"'").
						And("\"mck8sId\"", influxBuilder.Equal, "'"+info.ServiceID+"'").
						GroupByTime(timeCriteria).
						GroupByTag("\"nsId\"").
						GroupByTag("\"mck8sId\"").
						Fill("0").
						OrderByTime("ASC")
				}
				if strings.EqualFold(info.MCK8SReqInfo.GroupBy, types.Node) {
					query = query.Where("time", influxBuilder.MoreThan, timeDuration).
						And("\"nsId\"", influxBuilder.Equal, "'"+info.NsID+"'").
						And("\"mck8sId\"", influxBuilder.Equal, "'"+info.ServiceID+"'").
						And("\"node_name\"", influxBuilder.Equal, "'"+info.MCK8SReqInfo.Node+"'").
						GroupByTime(timeCriteria).
						GroupByTag("\"nsId\"").
						GroupByTag("\"mck8sId\"").
						GroupByTag("\"node_name\"").
						Fill("0").
						OrderByTime("ASC")
				}
			default:
				if strings.EqualFold(info.MCK8SReqInfo.GroupBy, types.Node) {
					query = query.Where("time", influxBuilder.MoreThan, timeDuration).
						And("\"nsId\"", influxBuilder.Equal, "'"+info.NsID+"'").
						And("\"mck8sId\"", influxBuilder.Equal, "'"+info.ServiceID+"'").
						And("\"node_name\"", influxBuilder.Equal, "'"+info.MCK8SReqInfo.Node+"'").
						GroupByTime(timeCriteria).
						GroupByTag("\"nsId\"").
						GroupByTag("\"mck8sId\"").
						GroupByTag("\"node_name\"").
						Fill("0").
						OrderByTime("ASC")
				}
				if strings.EqualFold(info.MCK8SReqInfo.GroupBy, types.Namespace) {
					query = query.Where("time", influxBuilder.MoreThan, timeDuration).
						And("\"nsId\"", influxBuilder.Equal, "'"+info.NsID+"'").
						And("\"mck8sId\"", influxBuilder.Equal, "'"+info.ServiceID+"'").
						And("\"namespace\"", influxBuilder.Equal, "'"+info.MCK8SReqInfo.Namespace+"'").
						GroupByTime(timeCriteria).
						GroupByTag("\"nsId\"").
						GroupByTag("\"mck8sId\"").
						GroupByTag("\"namespace\"").
						Fill("0").
						OrderByTime("ASC")
				}
				if strings.EqualFold(info.MCK8SReqInfo.GroupBy, string(types.MCK8S_POD)) {
					query = query.Where("time", influxBuilder.MoreThan, timeDuration).
						And("\"nsId\"", influxBuilder.Equal, "'"+info.NsID+"'").
						And("\"mck8sId\"", influxBuilder.Equal, "'"+info.ServiceID+"'").
						And("\"namespace\"", influxBuilder.Equal, "'"+info.MCK8SReqInfo.Namespace+"'").
						And("\"pod_name\"", influxBuilder.Equal, "'"+info.MCK8SReqInfo.Pod+"'").
						GroupByTime(timeCriteria).
						GroupByTag("\"nsId\"").
						GroupByTag("\"mck8sId\"").
						GroupByTag("\"namespace\"").
						GroupByTag("\"pod_name\"").
						Fill("0").
						OrderByTime("ASC")
				}
			}
		}
	} else {
		query = query.Where("time", influxBuilder.MoreThan, timeDuration).
			And("\"vmId\"", influxBuilder.Equal, "'"+info.VMID+"'").
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

func getPerSecMetric(info types.DBMetricRequestInfo, fieldArr ...string) string {
	var query string

	var timeCriteria string
	switch info.Period {
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
	if info.MonitoringMechanism {
		if util.CheckMCK8SType(info.ServiceType) {
			if strings.EqualFold(info.MCK8SReqInfo.GroupBy, types.Node) {
				whereQueryForm = " FROM \"%s\" WHERE time > (now()+1m) - %s AND \"nsId\"='%s' AND \"mck8sId\"='%s' AND \"node_name\"='%s' GROUP BY time(%s), \"nsId\", \"mck8sId\", \"node_name\" fill(0)"
				query += fmt.Sprintf(whereQueryForm, info.MetricName, info.Duration, info.NsID, info.ServiceID, info.MCK8SReqInfo.Node, timeCriteria)
			}
			if strings.EqualFold(info.MCK8SReqInfo.GroupBy, types.Namespace) {
				whereQueryForm = " FROM \"%s\" WHERE time > (now()+1m) - %s AND \"nsId\"='%s' AND \"mck8sId\"='%s' AND \"namespace\"='%s' GROUP BY time(%s), \"nsId\", \"mck8sId\", \"namespace\" fill(0)"
				query += fmt.Sprintf(whereQueryForm, info.MetricName, info.Duration, info.NsID, info.ServiceID, info.MCK8SReqInfo.Namespace, timeCriteria)
			}
			if strings.EqualFold(info.MCK8SReqInfo.GroupBy, string(types.MCK8S_POD)) {
				whereQueryForm = " FROM \"%s\" WHERE time > (now()+1m) - %s AND \"nsId\"='%s' AND \"mck8sId\"='%s' AND \"namespace\"='%s' AND \"pod_name\"='%s' GROUP BY time(%s), \"nsId\", \"mck8sId\", \"namespace\", \"pod_name\" fill(0)"
				query += fmt.Sprintf(whereQueryForm, info.MetricName, info.Duration, info.NsID, info.ServiceID, info.MCK8SReqInfo.Namespace, info.MCK8SReqInfo.Pod, timeCriteria)
			}
		} else {
			whereQueryForm = " FROM \"%s\" WHERE time > (now()+1m) - %s AND \"vmId\"='%s' GROUP BY time(%s) fill(0)"
			query += fmt.Sprintf(whereQueryForm, info.MetricName, info.Duration, info.VMID, timeCriteria)
		}
	} else {
		whereQueryForm = " FROM \"%s\" WHERE time > (now()+1m) - %s AND \"vmId\"='%s' GROUP BY time(%s), \"vmId\", \"nsId\", \"mcisId\" fill(0)"
		query += fmt.Sprintf(whereQueryForm, info.MetricName, info.Duration, info.VMID, timeCriteria)
	}

	return query
}
