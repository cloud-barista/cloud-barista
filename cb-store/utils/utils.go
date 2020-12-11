// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2019.08.

package utils

import (
	"strings"
	icbs "github.com/cloud-barista/cb-store/interfaces"
)

// GetNodeValue("/a/b/c", 2) => "b"
func GetNodeValue(strPath string, depth int) string {
        strList := strings.Split(strPath, "/")

        return strList[depth]
}

// CheckNodeValue("/a/b/c", 2, "b") => true
// CheckNodeValue("/a/b/c", 2, "c") => false
func CheckNodeValue(strPath string, depth int, word string) bool {
        strList := strings.Split(strPath, "/")

        if strList[depth] == word {
                return true
        }
        return false
}


// * key: /ns or /ns/
// * keyValueList:
// 	/ns/ns-01 => O
// 	/ns/ns-01/mcis/mcis-01 => X
func GetChildList(keyValueList []*icbs.KeyValue, key string) ([]*icbs.KeyValue) {
	// if key is "/ns" => "/ns/"
	key = strings.TrimSpace(key) + "/"       // => "/ns/" or "/ns//"
	key = strings.ReplaceAll(key, "//", "/") // => "/ns/"
	key = strings.ReplaceAll(key, "//", "/") // => "/ns/"

	var retKeyValueList []*icbs.KeyValue
	for _, kv := range keyValueList {

		// if erase "/ns/"
		// "/ns/ns-01" => "ns-01"  => if no "/", O
		// "/ns/ns-01/mcis/mcis-01" => "ns-01/mcis/mcis-01" => if has "/", X
		strList := strings.Split(kv.Key, key) // key = "/ns/"
		if len(strList) != 2 {
			continue;
		} else {
			is := strings.Index(strList[1], "/")
			if is == -1 { // strList[1] == "ns-01"
				// Got it!
				retKeyValueList = append(retKeyValueList, kv)		
			}
		}

	}
	return retKeyValueList
}
