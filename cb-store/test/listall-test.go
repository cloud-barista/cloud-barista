// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2019.08.

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

	cblog.Info("start test!!")

	// ## GetList
	keyValueList, _ := store.GetList("/", true) // true = Ascending

	fmt.Println("=========================== GetList(\"/\", Ascending)")
	for _, ev := range keyValueList {
		fmt.Println("<" + ev.Key + "> " + ev.Value)
	}
	fmt.Println("===========================")

	cblog.Info("finish test!!")
}
