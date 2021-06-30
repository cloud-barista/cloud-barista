// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2019.08.

package main

import (
        "fmt"

        "github.com/cloud-barista/cb-store/utils"
	icbs "github.com/cloud-barista/cb-store/interfaces"

	"testing"
)

func TestGetChildList_1(t *testing.T) {

	// ## Test Data & Specs
	keyValueData := []*icbs.KeyValue {
		{"/ns/ns-01", "ns-01-value"},
		{"/ns/ns-jhseo", "ns-jhseo-value"},
		{"/ns/ns-01/mcis/mcis-01", "mcis-01-value"},
	}

	keyValueList := utils.GetChildList(keyValueData, "/ns") 
	if len(keyValueList) != 1 {
		t.Errorf("The Number of chiid list is not one.(%d)", len(keyValueList))
	}

        fmt.Println("=================================")
        for _, ev := range keyValueList {
                fmt.Println("\t<" + ev.Key + "> " + ev.Value)
        }
        fmt.Println("=================================")
}

func TestGetChildList_2(t *testing.T) {

        // ## Test Data & Specs
        keyValueData := []*icbs.KeyValue {
                {"/ns/ns-01", "ns-01-value"},
                {"/ns/ns-01/mcis/mcis-01", "mcis-01-value"},
        }

        keyValueList := utils.GetChildList(keyValueData, "/ns/ns-01/mcis")
        if len(keyValueList) != 1 {
                t.Errorf("The Number of chiid list is not one.(%d)", len(keyValueList))
        }

        fmt.Println("=================================")
        for _, ev := range keyValueList {
                fmt.Println("\t<" + ev.Key + "> " + ev.Value)
        }
        fmt.Println("=================================")
}
