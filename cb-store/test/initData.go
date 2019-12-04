// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2019.08.

package main

import (
        "fmt"

	"github.com/sirupsen/logrus"
        "github.com/cloud-barista/cb-store"
	"github.com/cloud-barista/cb-store/config"
	icbs "github.com/cloud-barista/cb-store/interfaces"
)

var cblog *logrus.Logger
var store icbs.Store

func init() {
        cblog = config.Cblogger
        store = cbstore.GetStore()
}

func main() {

        cblog.Info("start test!!")
        fmt.Println("===========================initDB")

	// ## init data
	err := store.InitData()
	if err != nil {
		cblog.Error(err)
	}

        cblog.Info("finish test!!")
}
