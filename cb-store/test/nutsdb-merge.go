// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2022.06.

package main

import (
        "fmt"

        cbstore "github.com/cloud-barista/cb-store"
        "github.com/cloud-barista/cb-store/config"
        icbs "github.com/cloud-barista/cb-store/interfaces"
        "github.com/sirupsen/logrus"
)

var cblog *logrus.Logger
var store icbs.Store

func init() {
        cblog = config.Cblogger
        store = cbstore.GetStore()
}

func main() {

        cblog.Info("Start Merge....")

        // merge
        err := store.Merge()
                if nil != err {
                config.Cblogger.Error(err)
        }

        fmt.Println("===========================")

        cblog.Info("Finish Merge!!")
}
