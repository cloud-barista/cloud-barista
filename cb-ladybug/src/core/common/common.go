package common

import (
	cbstore "github.com/cloud-barista/cb-store"
	icbs "github.com/cloud-barista/cb-store/interfaces"
)

var CBStore icbs.Store

func init() {
	// cblog is a global variable.
	CBStore = cbstore.GetStore()
}
