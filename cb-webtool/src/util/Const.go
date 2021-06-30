package util

import (
	"os"
)

// 상수 정의

// var CloudConnectionUrl = os.Getenv("SPIDER_URL")
// var TumbleUrl = os.Getenv("TUMBLE_URL")
// var CloudConnectionUrl = os.Getenv("SPIDER_URL")
// var NameSpaceUrl = os.Getenv("TUMBLE_URL")

var SPIDER = os.Getenv("SPIDER_URL")
var TUMBLEBUG = os.Getenv("TUMBLE_URL")
var DRAGONFLY = os.Getenv("DRAGONFLY_URL")
var LADYBUG = os.Getenv("LADYBUG_URL")

var HTTP_CALL_SUCCESS = 200
var HTTP_POST_SUCCESS = 201

var PROVIDER_AWS = "aws"
var PROVIDER_AZURE = "azure"
var PROVIDER_ALIBABA = "alibaba"
var PROVIDER_GCP = "gcp"
var PROVIDER_CLOUDIT = "cloudit"
var PROVIDER_OPENSTACK = "openstack"
var PROVIDER_MOCK = "mock"

// MCIS의 상태(소문자).  (기타 상태는 UNDEFINED + ETC)
var MCIS_STATUS_RUNNING = "running"
var MCIS_STATUS_INCLUDE = "include"
var MCIS_STATUS_SUSPENDED = "suspended"
var MCIS_STATUS_TERMINATED = "terminated"
var MCIS_STATUS_PARTIAL = "partial"
var MCIS_STATUS_ETC = "etc"

// MCKA의 상태(소문자)
var MCKS_STATUS_RUNNING = "running"
var MCKS_STATUS_INCLUDE = "include"
var MCKS_STATUS_SUSPENDED = "suspended"
var MCKS_STATUS_TERMINATED = "terminated"
var MCKS_STATUS_PARTIAL = "partial"
var MCKS_STATUS_ETC = "etc"

var STATUS_ARRAY = []string{"running", "stopped", "terminated"}

// VM의 상태(소문자).  (기타 상태는 UNDEFINED + ETC)
var VM_STATUS_RUNNING = "running"

// var VM_STATUS_RESUMING = "Resuming"
var VM_STATUS_INCLUDE = "include"

var VM_STATUS_SUSPENDED = "suspended"
var VM_STATUS_STOPPED = "stopped"
var VM_STATUS_TERMINATED = "terminated"

// var VM_STATUS_UNDEFINED = "statusUndefined"
// var VM_STATUS_PARTIAL = "partial"
var VM_STATUS_ETC = "etc"

var VM_STATUS_ARRAY = []string{"running", "suspended", "terminated", "etc"}

/////// Store 저장 용 Const
var STORE_NAMESPACELIST = "namespacelist"
var STORE_CLOUDOSLIST = "cloudoslist"
var STORE_REGIONLIST = "regionlist"
var STORE_CREDENTIALLIST = "credentiallist"
var STORE_DRIVERLIST = "driverlist"
var STORE_CLOUDCONNECTIONCONFIGLIST = "cloudconnectionconfiglist"


var STORE_DEFAULT_NAMESPACENAME = "defaultnamespacename"
var STORE_DEFAULT_NAMESPACEID = "defaultnamespaceid"