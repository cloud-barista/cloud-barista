package metricstore

import (
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdbv1"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdbv2"
)

type StoreType string

const (
	InfluxDBV1Type StoreType = "influxdbv1"
	InfluxDBV2Type StoreType = "influxdbv2"
)

// Config is common interface of Monitoring Metric Storage Configuration
type Config interface{}

// Storage is common interface of Monitoring Metric Storage
type Storage interface {
	Init() error
	WriteMetric(metrics map[string]interface{}) error
	ReadMetric(vmId string, metric string, period string, aggregateType string, duration string) (interface{}, error)
}

// New Storage
func NewStorage(storeType StoreType, storageConfig Config) (Storage, error) {
	var sto Storage
	switch storeType {
	case InfluxDBV1Type:
		// InfluxDB v1.x
		if config, ok := storageConfig.(influxdbv1.Config); ok {
			sto = &influxdbv1.Storage{Config: config}
		} else {
			return nil, invalidConfigError(storeType)
		}
	case InfluxDBV2Type:
		// InfluxDB v2.x
		if conf, ok := storageConfig.(influxdbv2.Config); ok {
			sto = &influxdbv2.Storage{Config: conf}
		} else {
			return nil, invalidConfigError(storeType)
		}
	default:
		return nil, notSupportedTypeError(storeType)
	}
	// Storage Initialization
	if err := sto.Init(); err != nil {
		return nil, err
	}
	return sto, nil
}
