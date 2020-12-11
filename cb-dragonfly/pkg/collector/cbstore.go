package collector

import (
	"sync"

	cb "github.com/cloud-barista/cb-store"
	icbs "github.com/cloud-barista/cb-store/interfaces"
	"github.com/cloud-barista/cb-store/utils"
	"github.com/sirupsen/logrus"
)

type CBStore struct {
	Store icbs.Store
}

var once sync.Once
var cbstore CBStore

func Initialize() {
	cbstore.Store = cb.GetStore()
}

func GetInstance() *CBStore {
	once.Do(func() {
		Initialize()
	})
	return &cbstore
}

func (cs *CBStore) StorePut(key string, value string) {
	_ = cs.Store.Put(key, value)
}

func (cs *CBStore) StoreGet(key string) string {
	keyVal, _ := cs.Store.Get(key)
	if keyVal == nil {
		return ""
	}
	return keyVal.Value
}

func (cs *CBStore) StoreDelete(key string) {
	_ = cs.Store.Delete(key)
}

func (cs *CBStore) StoreGetList(key string, sortAscend bool) []string {
	keyVal, err := cs.Store.GetList(key, sortAscend)
	if err != nil {
		logrus.Debug(err)
		return []string{err.Error()}
	}
	result := []string{}
	for _, ev := range keyVal {
		if len(ev.Key) != 0 {
			result = append(result, ev.Key)
		}
	}
	return result
}

func (cs *CBStore) StoreDelList(key string) error {
	keyVal, err := cs.Store.GetList(key, true)
	if err != nil {
		logrus.Debug(err)
		return err
	}
	for _, ev := range keyVal {
		err = cs.Store.Delete(ev.Key)
		return err
	}
	return nil
}

func (cs *CBStore) StoreGetNodeValue(key string, depth int) string {
	return utils.GetNodeValue(key, depth)
}
