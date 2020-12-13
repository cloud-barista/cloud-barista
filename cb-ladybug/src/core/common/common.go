package common

import (
	cbstore "github.com/cloud-barista/cb-store"
	"github.com/cloud-barista/cb-store/config"
	icbs "github.com/cloud-barista/cb-store/interfaces"
	"github.com/sirupsen/logrus"
)

var CBLog *logrus.Logger
var CBStore icbs.Store

func init() {
	// cblog is a global variable.
	CBLog = config.Cblogger
	CBStore = cbstore.GetStore()
}
