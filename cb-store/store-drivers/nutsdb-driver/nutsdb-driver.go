// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2019.07.

package cbstore

import (
	"io/ioutil"
	"os"

	"github.com/cloud-barista/cb-store/config"
	"github.com/xujiajun/nutsdb"
	icbs "github.com/cloud-barista/cb-store/interfaces"
)

type NUTSDBDriver struct{}

var (
        db     *nutsdb.DB
        bucket string
)

func init() {
	initialize()
}

func initialize() {
        fileDir := config.GetConfigInfos().NUTSDB.DBPATH
        config.Cblogger.Info("######## dbfile: " + fileDir)

        opt := nutsdb.DefaultOptions
        opt.Dir = fileDir
        opt.SegmentSize = config.GetConfigInfos().NUTSDB.SEGMENTSIZE

        //opt.SegmentSize = 1024 * 1024 // 1MB
        db, _ = nutsdb.Open(opt)
        bucket = "bucketForString"
}

// If InitDB will be done online by other process, this is not effective until restart.
func (nutsdbDriver *NUTSDBDriver) InitDB() error {
	fileDir := config.GetConfigInfos().NUTSDB.DBPATH
        files, _ := ioutil.ReadDir(fileDir)
        for _, f := range files {
                name := f.Name()
                if name != "" {
                        err := os.RemoveAll(fileDir + "/" + name)
                        if err != nil {
				config.Cblogger.Error(err)
				return err
                        }
                }
        }
	return nil
}

// 1. get All data
// 2. delete All data
// If InitDB will be done online by other process, this is not effective until restart.
func (nutsdbDriver *NUTSDBDriver) InitData() error {
        config.Cblogger.Info("Call InitData")

	var keyList [][]byte

	// get all data
	if err := db.View(
	func(tx *nutsdb.Tx) error {
		entries, err := tx.GetAll(bucket)
		if err != nil {
			return err
		}
		keyList = make([][]byte, len(entries))
		for count, entry := range entries {
			keyList[count] = entry.Key
		}
		return nil
	}); err != nil {
		if err.Error() == "bucket is empty" {
			return nil
		}
                config.Cblogger.Error(err)
		return err
	}


	// delete all data
        if err := db.Update(
                func(tx *nutsdb.Tx) error {
		for _, one := range keyList {
			//key := []byte(one)
			if err := tx.Delete(bucket, one); err != nil {
				config.Cblogger.Error(err)
				return err
			}
		} // end of for
                return nil
        }); err != nil {
                config.Cblogger.Error(err)
                return err
        }
        return nil

}

func (nutsdbDriver *NUTSDBDriver) Put(key string, value string) error {
	config.Cblogger.Info("Key:" + key  + ", value:" + value)

        if err := db.Update(
                func(tx *nutsdb.Tx) error {
                        key := []byte(key)
                        val := []byte(value)
                        return tx.Put(bucket, key, val, 0)
                }); err != nil {
			config.Cblogger.Error(err)
			return err
        }

	return nil
}

func (nutsdbDriver *NUTSDBDriver) Get(key string) (*icbs.KeyValue, error) {
	config.Cblogger.Info("Key:" + key)

	var keyValue *icbs.KeyValue
        if err := db.View(
                func(tx *nutsdb.Tx) error {
                        key := []byte(key)
                        e, err := tx.Get(bucket, key)
                        if err != nil {
				if err.Error() == "key not found" {
					keyValue = nil
					return nil
				}
				config.Cblogger.Error(err)
                                return err
                        }
			keyValue = &icbs.KeyValue{string(key), string(e.Value)}
			return nil
                }); err != nil {
			config.Cblogger.Error(err)
			return nil, err
		}

	return keyValue, nil
}

func (nutsdbDriver *NUTSDBDriver) GetList(key string, sortAscend bool) ([]*icbs.KeyValue, error) {
        config.Cblogger.Info("Key:" + key)

        var keyValueList []*icbs.KeyValue
        if err := db.View(
                func(tx *nutsdb.Tx) error {
                        key := []byte(key)
                        entries, _, err := tx.PrefixScan(bucket, key, 0, 10000)
//config.Cblogger.Infof("================================:%v, %v", key, len(entries))
                        if err != nil {
				if err.Error() == "prefix scans not found" {
					return nil
				}
                                config.Cblogger.Error(err)
                                return err
                        }
			keyValueList = make([]*icbs.KeyValue, len(entries))
			if sortAscend {
				for k, entry := range entries {
					tmpOne := icbs.KeyValue{string(entry.Key), string(entry.Value)}
					keyValueList[k] = &tmpOne
				}
			} else {
				for k, entry := range entries {
					tmpOne := icbs.KeyValue{string(entry.Key), string(entry.Value)}
					keyValueList[len(entries)-1-k] = &tmpOne
				}
			}
                        return nil
                }); err != nil {
                        config.Cblogger.Error(err)
                        return nil, err
                }

        return keyValueList, nil
}

func (nutsdbDriver *NUTSDBDriver)Delete(key string) error {
	config.Cblogger.Info("Key:" + key)

	if err := db.Update(
		func(tx *nutsdb.Tx) error {
		key := []byte(key)
		if err := tx.Delete(bucket, key); err != nil {
			config.Cblogger.Error(err)
			return err
		}
		return nil
	}); err != nil {
		config.Cblogger.Error(err)
		return err
	}
	return nil
}

