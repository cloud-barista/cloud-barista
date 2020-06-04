package influxdbv1

import (
	"errors"
	"fmt"
	influxBuilder "github.com/Scalingo/go-utils/influx"
	influxdbClient "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
	"time"
)

type ClientOptions struct {
	URL      string
	Username string
	Password string
}

type Config struct {
	ClientOptions []ClientOptions
	Database      string
}

type Storage struct {
	Config  Config
	Clients []influxdbClient.Client
}

func (s *Storage) Init() error {
	for _, c := range s.Config.ClientOptions {
		client, err := influxdbClient.NewHTTPClient(influxdbClient.HTTPConfig{
			Addr:     c.URL,
			Username: c.Username,
			Password: c.Password,
		})
		if err != nil {
			logrus.Error(err)
			return err
		}
		if _, _, err := client.Ping(time.Millisecond * 100); err != nil {
			logrus.Error(err)
			return err
		}

		q := influxdbClient.Query{
			Command:  fmt.Sprintf("create database %s", "cbmon"),
			Database: "cbmon",
			//RetentionPolicy: "",
			//Precision:       "",
			//Chunked:         false,
			//ChunkSize:       0,
			//Parameters:      nil,
		}

		// ignore the error of existing database
		client.Query(q)

		s.Clients = append(s.Clients, client)
	}
	return nil
}

//func (s *Storage) WriteMetric(metrics types.Metrics) error {
func (s *Storage) WriteMetric(metrics map[string]interface{}) error {

	bp, err := s.parseMetric(metrics)
	if err != nil {
		logrus.Error("Failed to parse collector metrics to influxdb v1")
		return err
	}
	for _, influx := range s.Clients {
		if err := influx.Write(bp); err != nil {
			logrus.Error("Failed to write influxdb")
			return err
		}
	}
	return nil
}

//func (s *Storage) ReadMetric(vmId string, metric string, duration string) (interface{}, error) {
func (s *Storage) ReadMetric(vmId string, metric string, period string, aggregateType string, duration string) (interface{}, error) {

	influx := s.Clients[0]

	queryString, err := s.buildQuery(vmId, metric, period, aggregateType, duration)
	if err != nil {
		return nil, err
	}
	query := influxdbClient.NewQuery(queryString, s.Config.Database, "")
	res, _ := influx.Query(query)

	if res.Err != "" {
		return nil, errors.New(res.Err)
	}

	if len(res.Results) != 0 {
		if len(res.Results[0].Series) != 0 {
			return res.Results[0].Series[0], nil
		}
	}
	return nil, nil
}

func (s *Storage) parseMetric(metrics map[string]interface{}) (influxdbClient.BatchPoints, error) {

	bp, err := s.newBatchPoints()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	for hostId, v := range metrics {
		tagArr := map[string]string{}
		tagArr["hostId"] = hostId

		vToMap := v.(map[string]interface{})
		tagMapsInterface := vToMap["tag"].(map[string]interface{})

		tapMapsString := make(map[string]string)
		for k, v := range tagMapsInterface {
			tapMapsString[k] = v.(string)
		}

		tagArr["mcisId"] = tapMapsString["mcisId"]
		tagArr["osType"] = tapMapsString["osType"]

		delete(vToMap, "tag")

		for metricName, metric := range vToMap {
			metricPoint, err := influxdbClient.NewPoint(metricName, tagArr, metric.(map[string]interface{}), now)
			if err != nil {
				logrus.Error("Failed to create InfluxDB metric point: ", err)
				continue
			}
			bp.AddPoint(metricPoint)
		}
	}
	//spew.Dump(bp)
	return bp, nil
}

func (s *Storage) newBatchPoints() (influxdbClient.BatchPoints, error) {
	// TODO: implements
	return influxdbClient.NewBatchPoints(influxdbClient.BatchPointsConfig{
		Database: s.Config.Database,
	})
}

func (s *Storage) buildQuery(vmId string, metric string, period string, aggregateType string, duration string) (string, error) {

	// 통계 기준 설정
	if aggregateType == "avg" {
		aggregateType = "mean"
	}

	// 시간 범위 설정
	timeDuration := fmt.Sprintf("(now()+1m) - %s", duration)

	// 시간 단위 설정
	var timeCriteria time.Duration
	switch period {
	case "m":
		timeCriteria = time.Minute
	case "h":
		timeCriteria = time.Hour
	case "d":
		timeCriteria = time.Hour * 24
	}

	// InfluXDB 쿼리 생성
	var query influxBuilder.Query

	switch metric {

	case "cpu":

		query = influxBuilder.NewQuery().On(metric).
			Field("usage_utilization", aggregateType).
			Field("usage_system", aggregateType).
			Field("usage_idle", aggregateType).
			Field("usage_iowait", aggregateType).
			Field("usage_irq", aggregateType).
			Field("usage_softirq", aggregateType).
			Field("usage_user", aggregateType).
			Field("usage_nice", aggregateType).
			Field("usage_steal", aggregateType).
			Field("usage_guest", aggregateType).
			Field("usage_guest_nice", aggregateType)

	case "cpufreq":
		query = influxBuilder.NewQuery().On(metric).
			Field("cur_freq", aggregateType)

	case "net":

		fieldArr := []string{"bytes_recv", "bytes_sent", "packets_recv", "packets_sent", "err_in", "err_out", "drop_in", "drop_out"}
		query := s.getPerSecMetric(vmId, metric, period, fieldArr, duration)
		return query, nil

	case "mem":

		query = influxBuilder.NewQuery().On(metric).
			Field("used_percent", aggregateType).
			Field("total", aggregateType).
			Field("used", aggregateType).
			Field("free", aggregateType).
			Field("shared", aggregateType).
			Field("buffered", aggregateType).
			Field("cached", aggregateType)

	case "disk":

		query = influxBuilder.NewQuery().On(metric).
			Field("used_percent", aggregateType).
			Field("total", aggregateType).
			Field("used", aggregateType).
			Field("free", aggregateType)

	case "diskio":

		fieldArr := []string{"read_bytes", "write_bytes", "reads", "writes", "read_time", "write_time"}
		query := s.getPerSecMetric(vmId, metric, period, fieldArr, duration)
		return query, nil

	default:
		return "", errors.New("not found metric")
	}

	query = query.Where("time", influxBuilder.MoreThan, timeDuration).
		And("\"hostId\"", influxBuilder.Equal, "'"+vmId+"'").
		GroupByTime(timeCriteria).
		GroupByTag("\"hostId\"").
		Fill("0").
		OrderByTime("ASC")

	queryString := query.Build()

	return queryString, nil
}

func (s *Storage) getPerSecMetric(vmId string, metric string, period string, fieldArr []string, duration string) string {
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

	// 메트릭 조회 조건 쿼리 생성
	whereQueryForm := " FROM \"%s\" WHERE time > (now()+1m) - %s AND \"hostId\"='%s' GROUP BY time(%s) fill(0)"
	query += fmt.Sprintf(whereQueryForm, metric, duration, vmId, timeCriteria)

	return query
}
