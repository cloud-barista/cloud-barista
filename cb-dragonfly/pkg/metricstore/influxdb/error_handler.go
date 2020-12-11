package influxdb

import (
	"github.com/pkg/errors"

	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdb/influxdbv1"
	"github.com/cloud-barista/cb-dragonfly/pkg/metricstore/influxdb/influxdbv2"
)

func invalidConfigError(storeType StoreType) error {
	msg := "invalid metricstore config"
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
	return errors.Errorf("metricstore type '%s' not supported", storeType)
}
