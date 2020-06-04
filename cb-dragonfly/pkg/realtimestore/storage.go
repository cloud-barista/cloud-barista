package realtimestore

import (
	"github.com/cloud-barista/cb-dragonfly/pkg/realtimestore/etcd"
	"go.etcd.io/etcd/client"
)

type StoreType string

const (
	ETCDV2Type StoreType = "etcdv2"
)

// Config is common interface of Realtime Monitoring Metric Storage Configuration
type Config interface {
}

// Storage is common interface of Realtime Monitoring Metric Storage
type Storage interface {
	Init() error
	WriteMetric(key string, metric interface{}) error
	ReadMetric(key string) (*client.Node, error)
	DeleteMetric(key string) error
}

// New Storage
func NewStorage(storeType StoreType, storageConfig Config) (Storage, error) {
	var sto Storage
	switch storeType {
	case ETCDV2Type:
		if config, ok := storageConfig.(etcd.Config); ok {
			sto = &etcd.Storage{Config: config}
		} else {
			return nil, invalidConfigError(storeType, config)
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
