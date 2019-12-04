package influxdbv2

import (
	"context"
	influxdbClient "github.com/influxdata/influxdb-client-go"
	"github.com/sirupsen/logrus"
	"time"
)

type Client struct {
	URL      string
	Username string
	Password string
}

type Config struct {
	Clients []Client
}

type Storage struct {
	Config  Config
	Clients []*influxdbClient.Client
}

func (d *Storage) Init() error {
	for _, c := range d.Config.Clients {
		influx, err := influxdbClient.New(nil,
			influxdbClient.WithAddress(c.URL),
			influxdbClient.WithUserAndPass(c.Username, c.Password))
		if err != nil {
			logrus.Error(err)
			return err
		}
		if err := influx.Ping(context.Background()); err != nil {
			logrus.Error(err)
			return err
		}
		d.Clients = append(d.Clients, influx)
	}
	return nil
}

func (d *Storage) WriteMetric(metrics map[string]interface{}) error {
	// TODO: implements
	// example code
	myMetrics := []influxdbClient.Metric{
		influxdbClient.NewRowMetric(
			map[string]interface{}{"memory": 1000, "cpu": 0.93},
			"system-metrics",
			map[string]string{"hostname": "hal9000"},
			time.Date(2019, 8, 1, 5, 6, 7, 8, time.UTC)),
		influxdbClient.NewRowMetric(
			map[string]interface{}{"memory": 1000, "cpu": 0.93},
			"system-metrics",
			map[string]string{"hostname": "hal9000"},
			time.Date(2019, 8, 1, 5, 6, 7, 9, time.UTC)),
	}
	for _, influx := range d.Clients {
		if err := influx.Write(context.Background(), "bucket", "org", myMetrics...); err != nil {
			return err
		}
	}
	return nil
}

func (d *Storage) ReadMetric(vmId string, metric string, period string, aggregateType string, duration string) (interface{}, error) {
	// TODO: implements
	//metrics := types.Metrics{}
	return nil, nil
}
