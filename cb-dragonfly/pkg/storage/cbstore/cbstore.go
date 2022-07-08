package cbstore

import (
	"strconv"
	"sync"

	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	cb "github.com/cloud-barista/cb-store"
	icbs "github.com/cloud-barista/cb-store/interfaces"
	"github.com/cloud-barista/cb-store/utils"
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

func (cs *CBStore) StorePut(key string, value string) error {
	return cs.Store.Put(key, value)
}

func (cs *CBStore) StoreGet(key string) (*string, error) {
	keyVal, err := cs.Store.Get(key)
	if err != nil {
		return nil, err
	}
	if keyVal == nil {
		return nil, nil
	}
	return &keyVal.Value, nil
}

func (cs *CBStore) StoreGetToInt(key string) int {
	keyVal, _ := cs.Store.Get(key)
	if keyVal == nil {
		return -1
	}
	returnIntVal, _ := strconv.Atoi(keyVal.Value)
	return returnIntVal
}

func (cs *CBStore) StoreGetToString(key string) string {
	keyVal, _ := cs.Store.Get(key)
	if keyVal == nil {
		return ""
	}
	return keyVal.Value
}

func (cs *CBStore) StoreDelete(key string) error {
	return cs.Store.Delete(key)
}

func (cs *CBStore) StoreGetListMap(key string, sortAscend bool) (map[string]string, error) {
	keyVal, err := cs.Store.GetList(key, sortAscend)
	if err != nil {
		util.GetLogger().Error(err)
		return nil, err
	}
	result := map[string]string{}
	for _, ev := range keyVal {
		if len(ev.Key) != 0 {
			result[ev.Key] = ev.Value
		}
	}
	return result, nil
}

func (cs *CBStore) StoreGetListArray(key string, sortAscend bool) ([]string, error) {
	keyVal, err := cs.Store.GetList(key, sortAscend)
	if err != nil {
		util.GetLogger().Error(err)
		return nil, err
	}
	var result []string
	for _, ev := range keyVal {
		if len(ev.Key) != 0 {
			result = append(result, ev.Key)
		}
	}
	return result, nil
}

func (cs *CBStore) StoreDelList(key string) error {
	keyVal, err := cs.Store.GetList(key, true)
	if err != nil {
		util.GetLogger().Error(err)
		return err
	}
	for _, ev := range keyVal {
		err = cs.Store.Delete(ev.Key)
		if err != nil {
			util.GetLogger().Error(err)
			return err
		}
	}
	return nil
}

func (cs *CBStore) StoreGetNodeValue(key string, depth int) string {
	return utils.GetNodeValue(key, depth)
}
