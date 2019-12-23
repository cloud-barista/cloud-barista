package metricstore

import (
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdbv1"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdbv2"
	"github.com/pkg/errors"
)

func invalidConfigError(storeType StoreType) error {
	msg := "invalid storage config"
	switch storeType {
	case InfluxDBV1Type:
		return errors.Errorf("%s: %v", msg, influxdbv1.Config{})
	case InfluxDBV2Type:
		return errors.Errorf("%s: %v", msg, influxdbv2.Config{})
	default:
		return errors.New(msg)
	}
}

func notSupportedTypeError(storeType StoreType) error {
	return errors.Errorf("storage type '%s' not supported", storeType)
}
