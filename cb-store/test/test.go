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
	"github.com/cloud-barista/cb-store/utils"
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
		{"/key1", "value"},
		{"/key1", "value1"},    // same key
		{"/key1/", "value2"},   // end with '/'
		{"/key1/%", "value3%"}, // with special char.
		{"/key1/key2/key3", "value4"},
		{"/space key", "space value5"},
		{"/newline \n key", "newline \n value6"},
		{"/a/b/c/123/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u", "value/value/value"},
	}

	// ### Put
	fmt.Println("=========================== Put(...)")
	for _, ev := range keyValueData {
		fmt.Println("<" + ev.Key + "> " + ev.Value)
		err := store.Put(ev.Key, ev.Value)
		if err != nil {
			cblog.Error(err)
		}
	}
	fmt.Println("===========================")

	// ## Get
	keyValue, _ := store.Get("/") // exact match

	fmt.Println("=========================== Get(\"/\")")
	fmt.Println("<" + keyValue.Key + "> " + keyValue.Value)
	fmt.Println("===========================")

	keyValue, _ = store.Get("/space key") // exact match

	fmt.Println("=========================== Get(\"space key\")")
	fmt.Println("<" + keyValue.Key + "> " + keyValue.Value)
	fmt.Println("===========================")

	// ## GetList
	keyValueList, _ := store.GetList("/", true) // true = Ascending

	fmt.Println("=========================== GetList(\"/\", Ascending)")
	for _, ev := range keyValueList {
		fmt.Println("<" + ev.Key + "> " + ev.Value)
	}
	fmt.Println("===========================")

	// ## GetList
	keyValueList, _ = store.GetList("/", false) // false = Descending

	fmt.Println("=========================== GetList(\"/\", Descending)")
	for _, ev := range keyValueList {
		fmt.Println("<" + ev.Key + "> " + ev.Value)
	}
	fmt.Println("===========================")

	// ## CheckNodeValue
	ret := utils.CheckNodeValue("/key1/key2/key3", 2, "key2")

	fmt.Println("===========================  utils.CheckNodeValue(\"/key1/key2/key3\", 2, \"key2\")")
	fmt.Printf("<exist?> %v\n", ret)

	ret2 := utils.CheckNodeValue("/key1/key2/key3", 2, "aws")

	fmt.Println("===========================  utils.CheckNodeValue(\"/key1/key2/key3\", 2, \"aws\")")
	fmt.Printf("<exist?> %v\n", ret2)

	fmt.Println("===========================")

	// ## Delete
	for _, ev := range keyValueData {
		err := store.Delete(ev.Key)
		if err != nil {
			cblog.Error(err)
		}
	}

	// ## GetNodeValue
	value := utils.GetNodeValue("/key1/key2/key3", 2)

	fmt.Println("===========================  utils.GetNodeValue(\"/key1/key2/key3\", 2)")
	fmt.Println("<Value> " + value)
	fmt.Println("===========================")

	// ## Delete
	for _, ev := range keyValueData {
		err := store.Delete(ev.Key)
		if err != nil {
			cblog.Error(err)
		}
	}

	// ## empty result
	keyValue, _ = store.Get("/powerkim/") // no key
	fmt.Println("=========================== Get(): no key")
	if keyValue != nil {
		fmt.Println("<" + keyValue.Key + "> " + keyValue.Value)
	}

	keyValueList, _ = store.GetList("/powerkim/", true) // no key list
	fmt.Println("=========================== GetList(): no key")
	for _, ev := range keyValueList {
		fmt.Println("<" + ev.Key + "> " + ev.Value)
	}
	fmt.Println("===========================")

	cblog.Info("finish test!!")
}
