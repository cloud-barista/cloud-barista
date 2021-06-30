package influxdb

import (
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdb/v1"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdb/v2"
	"github.com/pkg/errors"
)

type StoreType string

const (
	V1 StoreType = "v1"
	V2 StoreType = "v2"
)

// Config is common interface of Monitoring Metric Storage Configuration
type Config interface{}

// Storage is common interface of Monitoring Metric Storage
type Storage interface {
	Initialize() error
	WriteMetric(database string, metrics map[string]interface{}) error
	ReadMetric(isPush bool, nsId string, mcisId string, vmId string, metric string, period string, function string, duration string) (interface{}, error)
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
	case V2:
		// InfluxDB v2
		if conf, ok := storageConfig.(v2.Config); ok {
			storage = &v2.Storage{Config: conf}
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
	case V2:
		return errors.Errorf("%s: %v", msg, v2.Config{})
	default:
		return errors.New(msg)
	}
}
