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
        _ "github.com/cloud-barista/cb-store/utils"
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

	keyValue, err := store.Get("/no") // exact match

if keyValue == nil {
	fmt.Println("Get(): nil")
	fmt.Printf("Get() error: %v\n", err)
}else {
	
	fmt.Println("=========================== Get(\"/no\")")
	fmt.Println("<" + keyValue.Key + "> " + keyValue.Value)
	fmt.Println("===========================")
}

        // ## GetList
        keyValueList, err2 := store.GetList("/no", true) // true = Ascending
	fmt.Printf("store.GetList() error: %v\n", err2)

        fmt.Println("=========================== GetList(\"/no\", Ascending)")
        for _, ev := range keyValueList {
                fmt.Println("<" + ev.Key + "> " + ev.Value)
        }
        fmt.Println("===========================")

        cblog.Info("finish test!!")
}
