// Package influxdb - Metrics 처리를 위한 InfluxDB Client 기능 제공
package influxdb

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics/influxdb/counter"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics/influxdb/gauge"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/middlewares/metrics/influxdb/histogram"
	"github.com/influxdata/influxdb/client/v2"
)

// ===== [ Constants and Variables ] =====

var (
	errNoConfig = errors.New("no config present for the influxdb related with metircs middleware")
)

// ===== [ Types ] =====

// Config - InfluxDB 운영을 위한 설정 구조 (Metrics 와의 Cycle Import 문제로 재 정의)
type Config struct {
	Address         string        `yaml:"address"`
	ReportingPeriod time.Duration `yaml:"reporting_period"`
	BufferSize      int           `yaml:"buffer_size"`
	UserName        string        `yaml:"user_name"`
	Password        string        `yaml:"password"`
	Database        string        `yaml:"database"`
}

// clientWrapper - InfluxDB client 와 Metrics 정보를 관리하는 구조 정의
type clientWrapper struct {
	influxClient client.Client
	stats        func() interface{}
	logger       *logging.Logger
	db           string
	buff         *Buffer
}

// ===== [ Implementations ] =====

// keepUpdated - 주기적 (ReportingPeriod)으로 Metrics 데이터를 InfluxDB에 반영
func (cw clientWrapper) keepUpdated(ctx context.Context, ticker <-chan time.Time) {
	hostname, err := os.Hostname()
	if err != nil {
		cw.logger.Error("influxdb resolving the local hostname:", err.Error())
	}

	for {
		select {
		case <-ticker:
		case <-ctx.Done():
			return
		}

		// Collection Function을 호출해서 Metric에 수집되어 있는 Stats에 대한 Snapshot 추출
		snapshot := cw.stats().(metrics.Stats)

		if shouldSendPoints := len(snapshot.Counters) > 0 || len(snapshot.Gauges) > 0; !shouldSendPoints {
			continue
		}

		bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  cw.db,
			Precision: "s",
		})
		now := time.Unix(0, snapshot.Time)

		for _, p := range counter.Points(hostname, now, snapshot.Counters, *cw.logger) {
			bp.AddPoint(p)
		}

		for _, p := range gauge.Points(hostname, now, snapshot.Gauges, *cw.logger) {
			bp.AddPoint(p)
		}

		for _, p := range histogram.Points(hostname, now, snapshot.Histograms, *cw.logger) {
			bp.AddPoint(p)
		}

		if err := cw.influxClient.Write(bp); err != nil {
			cw.logger.Error("writing to influx:", err.Error())
			cw.buff.Add(bp)
			continue
		}

		pts := []*client.Point{}
		bpPending := cw.buff.Elements()
		for _, failedBP := range bpPending {
			pts = append(pts, failedBP.Points()...)
		}

		retryBatch, _ := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  cw.db,
			Precision: "s",
		})
		retryBatch.AddPoints(pts)

		if err := cw.influxClient.Write(retryBatch); err != nil {
			cw.logger.Error("writting to influx:", err.Error())
			cw.buff.Add(bpPending...)
			continue
		}
	}
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// SetupAndRun - InfluxDB로 Metric 처리를 반영하기 위한 Client 생성하고 주기적으로 Metrics 의 수집된 Stats 처리
func SetupAndRun(ctx context.Context, idbConf Config, collectFunc func() interface{}, logger *logging.Logger) error {
	// creating a new influxdb client

	if &idbConf == nil || idbConf.Address == "" {
		return errNoConfig
	}

	// InfluxDB client 설정 생성
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     idbConf.Address,
		Username: idbConf.UserName,
		Password: idbConf.Password,
		Timeout:  10 * time.Second,
	})
	if err != nil {
		return err
	}

	// InfluxDB Ping 테스트
	go func() {
		_, _, err := influxClient.Ping(time.Second)
		if err != nil {
			logger.Error("unable to ping the influxdb server:", err.Error())
			return
		}
	}()

	// 설정된 ReportingPeriod 기준으로 InfluxDB tick 생성
	tick := time.NewTicker(idbConf.ReportingPeriod)

	cw := clientWrapper{
		influxClient: influxClient,
		stats:        collectFunc,
		logger:       logger,
		db:           idbConf.Database,
		buff:         NewBuffer(idbConf.BufferSize),
	}

	// 지정된 주기단위로 InfluxDB에 Metrics 정보 처리
	go cw.keepUpdated(ctx, tick.C)

	return nil
}
