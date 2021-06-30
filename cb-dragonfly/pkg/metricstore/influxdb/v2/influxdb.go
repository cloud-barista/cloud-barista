package v2

import (
	influxdbClient "github.com/influxdata/influxdb-client-go"
)

type Client struct {
	URL      string
	Username string
	Password string
}

type Config struct {
	Client Client
}

type Storage struct {
	Config Config
	Client influxdbClient.Client
}

func (s Storage) Initialize() error {
	// TODO: implements
	return nil
}

func (s Storage) WriteMetric(database string, metrics map[string]interface{}) error {
	// TODO: implements
	return nil
}

func (s Storage) ReadMetric(isPush bool, nsId string, mcisId string, vmId string, metric string, period string, function string, duration string) (interface{}, error) {
	// TODO: implements
	return nil, nil
}

func (s Storage) DeleteMetric(database string, metric string, duration string) error {
	// TODO: implements
	return nil
}
