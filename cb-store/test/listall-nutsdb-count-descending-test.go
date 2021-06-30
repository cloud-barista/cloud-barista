// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2019.08.

package main

import (
	"fmt"
	"strconv"
	"time"

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

        // ## Test Data & Specs
        keyValueData := []icbs.KeyValue{
                {"/", "root"},
                {"/a", "value"},
                {"/c", "value"},
                {"/b", "value"},
                {"/3", "value"},
                {"/d", "value"},
                {"/a/a", "value"},
                {"/a/b", "value"},
                {"/b/b", "value"},
                {"/b/c", "value"},
        }


        // ### Put
        fmt.Println("=========================== Put(...)")
	startTime := time.Now()
        for _, ev := range keyValueData {
                //fmt.Println("<" + ev.Key + "> " + ev.Value)
                err := store.Put(ev.Key, ev.Value)
                if err != nil {
                        cblog.Error(err)
                }
        }
	elapsed := fmt.Sprintf("%.4f\n", time.Since(startTime).Seconds())
	fmt.Println("=========================== Put Key #: " + strconv.Itoa(len(keyValueData)) + ", elapsed time(sec): " + elapsed )
        fmt.Println("")
        fmt.Println("")


        // ### Put more 10,000
        fmt.Println("=========================== Put(...) over 10,000")
	keyNum := 100000
	startTime = time.Now()
	for i :=0; i< keyNum; i++ {
                err := store.Put(keyValueData[len(keyValueData)-1].Key + "/" + strconv.Itoa(i), keyValueData[len(keyValueData)-1].Value)
                if err != nil {
                        cblog.Error(err)
                }
        }
	elapsed = fmt.Sprintf("%.4f\n", time.Since(startTime).Seconds())

	fmt.Println("=========================== Put Key #: " + strconv.Itoa(keyNum) + ", elapsed time(sec): " + elapsed )
        fmt.Println("")
        fmt.Println("")

	startTime = time.Now()
	// ## GetList ascending
	keyValueList, _ := store.GetList("/", true) // true = Ascending
	elapsed = fmt.Sprintf("%.4f\n", time.Since(startTime).Seconds())

	fmt.Println("=========================== GetList(\"/\", Ascending): Key #: " + strconv.Itoa(len(keyValueList)) + ", elapsed time(sec): " + elapsed)
	/*
	for _, ev := range keyValueList {
		fmt.Println("<" + ev.Key + "> " + ev.Value)
	}
	*/

        fmt.Println("")
        fmt.Println("")


	startTime = time.Now()
        // ## GetList descending
        keyValueList, _ = store.GetList("/", false) // true = Ascending
	elapsed = fmt.Sprintf("%.4f\n", time.Since(startTime).Seconds())

	fmt.Println("=========================== GetList(\"/\", Descending): Key #: " + strconv.Itoa(len(keyValueList)) + ", elapsed time(sec): " + elapsed)
	/*
        for _, ev := range keyValueList {
                fmt.Println("<" + ev.Key + "> " + ev.Value)
        }
	*/
        fmt.Println("")
        fmt.Println("")

	cblog.Info("finish test!!")
}
