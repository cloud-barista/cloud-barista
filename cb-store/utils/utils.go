// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2019.08.

package utils

import (
	"strings"
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

