package v1

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	influxdbClient "github.com/influxdata/influxdb1-client/v2"
)

const (
	DefaultDatabase       = "cbmon"
	PullDatabase          = "cbmonpull"
	CBRetentionPolicyName = "df_rp"
)

type Config struct {
	Addr     string
	Username string
	Password string
	Database string
}

type Storage struct {
	Config Config
	Client influxdbClient.Client
}

var once sync.Once
var storage Storage

func GetInstance() *Storage {
	once.Do(func() {
		if storage.Client == nil {
			storage.Initialize()
		}
	})
	return &storage
}

func (s Storage) Initialize() error {
	if s.Config == (Config{}) {
		influxDBConfig := config.GetInstance().GetInfluxDBConfig()
		influxDBAddr := fmt.Sprintf("%s:%d", influxDBConfig.EndpointUrl, influxDBConfig.ExternalPort)
		s.Config.Addr = influxDBAddr
		s.Config.Username = influxDBConfig.UserName
		s.Config.Password = influxDBConfig.Password
	}

	client, err := influxdbClient.NewHTTPClient(influxdbClient.HTTPConfig{
		Addr:     s.Config.Addr,
		Username: s.Config.Username,
		Password: s.Config.Password,
	})
	if err != nil {
		util.GetLogger().Error(fmt.Sprintf("failed to create InfluxDB HTTP Client, error=%s", err))
		return err
	}
	if _, _, err := client.Ping(5 * time.Second); err != nil {
		util.GetLogger().Error(fmt.Sprintf("failed to ping InfluxDB, error=%s", err))
		return err
	}

	q1 := influxdbClient.Query{
		Command: fmt.Sprintf("create database %s", DefaultDatabase),
	}
	// ignore the error of existing database
	client.Query(q1)

	q2 := influxdbClient.Query{
		Command: fmt.Sprintf("create database %s", PullDatabase),
	}
	// ignore the error of existing database
	client.Query(q2)

	// cbmon rp 조회 후 없을 시 rp 생성
	if isRPonCBMonExist := s.checkDBRetionPolicy(client, DefaultDatabase); !isRPonCBMonExist {
		createRPq1 := influxdbClient.Query{
			Command: fmt.Sprintf("create retention policy %s on %s duration %s replication 1 default", CBRetentionPolicyName, DefaultDatabase, config.GetInstance().InfluxDB.RetentionPolicyDuration),
		}
		// influxdb rpDuration, shardGroupDuration 특성으로 인한 에러 검출
		_, err := client.Query(createRPq1)
		if err != nil {
			return err
		}
	}

	// cbmonpull rp 조회 후 없을 시 rp 생성
	if isRPonCBMonPullExist := s.checkDBRetionPolicy(client, PullDatabase); !isRPonCBMonPullExist {
		createRPq2 := influxdbClient.Query{
			Command: fmt.Sprintf("create retention policy %s on %s duration %s replication 1 default", CBRetentionPolicyName, PullDatabase, config.GetInstance().InfluxDB.RetentionPolicyDuration),
		}

		// influxdb rpDuration, shardGroupDuration 특성으로 인한 에러 검출
		_, err := client.Query(createRPq2)
		if err != nil {
			return err
		}
	}

	s.Client = client
	storage = s
	return nil
}

func (s Storage) checkDBRetionPolicy(client influxdbClient.Client, dbName string) bool {
	// Retention Policy 조회
	listQuery := influxdbClient.Query{
		Command: fmt.Sprintf("show retention policies on %s", dbName),
	}

	res, _ := client.Query(listQuery)

	isRPExist := false

	for _, db := range res.Results {
		for _, v := range db.Series {
			for _, rp := range v.Values {
				for _, targetValue := range rp {
					if targetValue == CBRetentionPolicyName {
						isRPExist = true
						break
					}
				}
			}
		}
	}
	return isRPExist
}

func (s Storage) WriteMetric(dbName string, metrics map[string]interface{}) error {
	bp, err := influxdbClient.NewBatchPoints(influxdbClient.BatchPointsConfig{
		Database: dbName,
	})
	if err != nil {
		return err
	}

	var tagInfo map[string]string
	tagInfo = nil
	now := time.Now().UTC()

	for _, metricVal := range metrics {
		metricValMap := metricVal.(map[string]interface{})
		if tagInfo == nil {
			tagInfo = metricValMap["tagInfo"].(map[string]string)
		}
		delete(metricValMap, "tagInfo")
		for metricName, metric := range metricValMap {
			metricPoint, err := influxdbClient.NewPoint(metricName, tagInfo, metric.(map[string]interface{}), now)
			if err != nil {
				util.GetLogger().Error("failed to create InfluxDB metricVal point: ", err)
				continue
			}
			bp.AddPoint(metricPoint)
		}
	}

	if err := s.Client.Write(bp); err != nil {
		util.GetLogger().Error("failed to write InfluxDB")
		return err
	}
	return nil
}

func (s Storage) WriteOnDemandMetric(dbName string, metricName string, tagArr map[string]string, metricVal map[string]interface{}) error {
	bp, err := influxdbClient.NewBatchPoints(influxdbClient.BatchPointsConfig{
		Database: dbName,
	})
	if err != nil {
		util.GetLogger().Error("failed to create InfluxDB metric point: ", err)
		return err
	}

	now := time.Now().UTC()
	metricPoint, err := influxdbClient.NewPoint(metricName, tagArr, metricVal, now)
	if err != nil {
		util.GetLogger().Error("failed to create InfluxDB metric point: ", err)
		return err
	}
	bp.AddPoint(metricPoint)

	if err := s.Client.Write(bp); err != nil {
		util.GetLogger().Error("failed to write InfluxDB")
		return err
	}
	return nil
}

func (s Storage) ReadMetric(isPush bool, nsId string, mcisId string, vmId string, metric string, period string, function string, duration string) (interface{}, error) {
	var database string
	if isPush {
		database = DefaultDatabase
	} else {
		database = PullDatabase
	}
	queryString, err := BuildQuery(isPush, vmId, metric, period, function, duration)
	if err != nil {
		return nil, err
	}
	query := influxdbClient.NewQuery(queryString, database, "")
	res, _ := s.Client.Query(query)
	if res.Err != "" {
		return nil, errors.New(res.Err)
	}
	if len(res.Results) > 0 {
		if len(res.Results[0].Series) > 0 {
			return res.Results[0].Series[0], nil
		}
	}
	return nil, nil
}

func (s Storage) DeleteMetric(database string, metric, duration string) error {
	whereQuery := "DELETE FROM \"%s\" WHERE time < now() + 1m - %s"
	query := influxdbClient.NewQuery(fmt.Sprintf(whereQuery, metric, duration), database, "")
	res, _ := s.Client.Query(query)
	if res.Err != "" {
		return errors.New(res.Err)
	}
	return nil
}
