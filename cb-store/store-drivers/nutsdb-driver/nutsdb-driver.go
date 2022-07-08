// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2019.07.

package cbstore

import (
	"io/ioutil"
	_ "fmt"
	"os"
	"strings"

	"github.com/cloud-barista/cb-store/config"
	icbs "github.com/cloud-barista/cb-store/interfaces"
	"github.com/xujiajun/nutsdb"
)

// NUTSDBDriver - NutsDB 처리 정보 구조
type NUTSDBDriver struct{}

var (
	db     *nutsdb.DB
	bucket string
)

func init() {
	// 패키지 로그 검증용
	//initialize()
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

// InitDB - Database 초기화
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

// InitData - 데이터 초기화
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

// Put - 지정한 키/값을 NutsDB 데이터로 추가
func (nutsdbDriver *NUTSDBDriver) Put(key string, value string) error {
	config.Cblogger.Info("Key:" + key + ", value:" + value)

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

// Get - 지정한 키의 값을 NutsDB 데이터에서 추출
func (nutsdbDriver *NUTSDBDriver) Get(key string) (*icbs.KeyValue, error) {
	config.Cblogger.Info("Key:" + key)

	var keyValue *icbs.KeyValue
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			key := []byte(key)
			e, err := tx.Get(bucket, key)
			if err != nil {
				if strings.Contains(err.Error(), "key not found") || strings.Contains(err.Error(), "not found bucket:") {
					keyValue = nil
					return nil
				}
				config.Cblogger.Error(err)
				return err
			}
			keyValue = &icbs.KeyValue{Key: string(key), Value: string(e.Value)}
			return nil
		}); err != nil {
		config.Cblogger.Error(err)
		return nil, err
	}

	return keyValue, nil
}

// GetList - 지정한 키와 정렬 조건을 기준으로 NustDB 데이터에서 추출
func (nutsdbDriver *NUTSDBDriver) GetList(key string, sortAscend bool) ([]*icbs.KeyValue, error) {
	config.Cblogger.Info("Key:" + key)

	keyValueList := []*icbs.KeyValue{}
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			key := []byte(key)
			offsetNum := 0
			limitNum := 10000 
			for true {
				entries, _, err := tx.PrefixScan(bucket, key, offsetNum, limitNum)
				//fmt.Printf("\n================================:%s, %d\n", key, len(entries))
				//config.Cblogger.Infof("================================:%v, %v", key, len(entries))
				if err != nil {
					if err.Error() == "prefix scans not found" {
						return nil
					}
					config.Cblogger.Error(err)
					return err
				}
				if len(entries) == 0 {
					return nil;
				}

				for _, entry := range entries {
					tmpOne := icbs.KeyValue{Key: string(entry.Key), Value: string(entry.Value)}
					keyValueList = append(keyValueList, &tmpOne)
				}

			offsetNum = offsetNum + limitNum	
			}  // end of for true

			return nil
		}); err != nil {
		config.Cblogger.Error(err)
		return nil, err
	}

	// for descending order
	if !sortAscend {
		len := len(keyValueList)
		for i := 0; i<len/2; i++ {  // swap the list around the center.
			keyValueList[i], keyValueList[len-1-i] = keyValueList[len-1-i], keyValueList[i]
		}
	}

	return keyValueList, nil
}

// Delete - 지정한 키의 데이터를 NutsDB 데이터에서 삭제
func (nutsdbDriver *NUTSDBDriver) Delete(key string) error {
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

// Close - NutsDB 연결 해제
func (nutsdbDriver *NUTSDBDriver) Close() error {
	err := db.Close()
	if nil != err {
		config.Cblogger.Error(err)
	}
	return err
}

// InitializeDriver - NutsDB Driver 초기화
func InitializeDriver() {
	initialize()
}

func (nutsdbDriver *NUTSDBDriver) Merge() error {
        err := db.Merge()
        if nil != err {
                config.Cblogger.Error(err)
        }
        return err
}
