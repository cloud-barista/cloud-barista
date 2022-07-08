package influxdb

import (
	"github.com/cloud-barista/cb-dragonfly/pkg/storage/metricstore/influxdb/v1"
	"github.com/cloud-barista/cb-dragonfly/pkg/types"
	"github.com/pkg/errors"
)

type StoreType string

const (
	V1 StoreType = "v1"
)

// Config is common interface of Monitoring Metric Storage Configuration
type Config interface{}

// Storage is common interface of Monitoring Metric Storage
type Storage interface {
	Initialize() error
	WriteMetric(database string, metrics map[string]interface{}) error
	ReadMetric(info types.DBMetricRequestInfo) (interface{}, error)
	DeleteMetric(database string, metric string, duration string) error
}

// NewStorage fuction for initialize InfluxDB Client
func NewStorage(storeType StoreType, storageConfig Config) error {
	var storage Storage
	switch storeType {
	case V1:
		// InfluxDB v1
		if config, ok := storageConfig.(v1.Config); ok {
			storage = v1.Storage{Config: config}
		} else {
			return invalidConfigError(storeType)
		}
	default:
		return errors.Errorf("InfluxDB %s not supported", storeType)
	}
	if err := storage.Initialize(); err != nil {
		return err
	}
	return nil
}

func invalidConfigError(storeType StoreType) error {
	msg := "invalid configuration of influxDB"
	switch storeType {
	case V1:
		return errors.Errorf("%s: %v", msg, v1.Config{})
	default:
		return errors.New(msg)
	}
}
