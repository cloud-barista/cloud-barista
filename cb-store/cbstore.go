// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2019.08.

package cbstore

import (
	"fmt"

	"github.com/cloud-barista/cb-store/config"
	icbs "github.com/cloud-barista/cb-store/interfaces"
	etcddrv "github.com/cloud-barista/cb-store/store-drivers/etcd-driver"
	nutsdrv "github.com/cloud-barista/cb-store/store-drivers/nutsdb-driver"
)

var configInfo *config.CBSTORECONFIG

func init() {
	//config.Cblogger.Info("calling!")
	configInfo = config.GetConfigInfos()
}

// initialize db
func InitStore() error {
	// NUTSDB: If InitDB will be done online by other process, this is not effective until restart.
	if configInfo.STORETYPE == "NUTSDB" {
		// 1. remove path: rm -rf ./meta_store/*
		// init nutsdb metainfo
		store := nutsdrv.NUTSDBDriver{}
		return store.InitDB()
	}
	if configInfo.STORETYPE == "ETCD" {
		// init etcd metainfo
		store := etcddrv.ETCDDriver{}
		return store.InitDB()
	}
	return fmt.Errorf("STORETYPE:" + configInfo.STORETYPE + " is not supported!!")
}

// clean all
func InitData() error {
	// NUTSDB: If InitDB will be done online by other process, this is not effective until restart.
	if configInfo.STORETYPE == "NUTSDB" {
		// 1. remove path: rm -rf ./meta_store/*
		// init nutsdb metainfo
		store := nutsdrv.NUTSDBDriver{}
		return store.InitData()
	}
	if configInfo.STORETYPE == "ETCD" {
		// init etcd metainfo
		store := etcddrv.ETCDDriver{}
		return store.InitData()
	}
	return fmt.Errorf("STORETYPE:" + configInfo.STORETYPE + " is not supported!!")
}

func GetStore() icbs.Store {
	if configInfo.STORETYPE == "NUTSDB" {
		// NutsDB 드라이버 초기화 호출
		nutsdrv.InitializeDriver()
		store := nutsdrv.NUTSDBDriver{}
		return &store
	}
	if configInfo.STORETYPE == "ETCD" {
		// ETCD 드라이버 초기화 호출
		etcddrv.InitializeDriver()
		store := etcddrv.ETCDDriver{}
		return &store
	}
	config.Cblogger.Errorf("STORETYPE:" + configInfo.STORETYPE + " is not supported!!")

	return nil
}
