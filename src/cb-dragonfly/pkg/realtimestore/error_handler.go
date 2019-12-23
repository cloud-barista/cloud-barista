package realtimestore

import (
	"github.com/cloud-barista/cb-dragonfly/pkg/realtimestore/etcd"
	"github.com/pkg/errors"
)

func invalidConfigError(storeType StoreType, config Config) error {
	msg := "invalid storage config"
	switch storeType {
	case ETCDV2Type:
		return errors.Errorf("%s %v", msg, config.(etcd.Config))
	default:
		return errors.New(msg)
	}
}

func notSupportedTypeError(storeType StoreType) error {
	return errors.Errorf("storage type '%s' not supported", storeType)
}
