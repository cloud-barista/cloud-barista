// Rest Runtime Server of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by CB-Spider Team, 2019.10.

package restruntime

import (
	"fmt"
	"time"

	"os"
	"net/http"
	"github.com/chyeh/pubip"

	cr "github.com/cloud-barista/cb-spider/api-runtime/common-runtime"
	aw "github.com/cloud-barista/cb-spider/api-runtime/rest-runtime/admin-web"
	"github.com/cloud-barista/cb-store/config"
	"github.com/sirupsen/logrus"

	// REST API (echo)
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var cblog *logrus.Logger

func init() {
	cblog = config.Cblogger
	currentTime := time.Now()
	cr.StartTime = currentTime.Format("2006.01.02 15:04:05 Mon")
	cr.MiddleStartTime = currentTime.Format("2006.01.02.15:04:05")
	cr.ShortStartTime = fmt.Sprintf("T%02d:%02d:%02d", currentTime.Hour(), currentTime.Minute(), currentTime.Second())
	cr.HostIPorName = getHostIPorName()
}

// REST API Return struct for boolena type
type BooleanInfo struct {
	Result string // true or false
}

type StatusInfo struct {
	Status string // PENDING | RUNNING | SUSPENDING | SUSPENDED | REBOOTING | TERMINATING | TERMINATED
}

//ex) {"POST", "/driver", registerCloudDriver}
type route struct {
	method, path string
	function     echo.HandlerFunc
}

func getHostIPorName() string {
	if os.Getenv("LOCALHOST") ==  "ON" {
		return "localhost"
	}

	ip, err := pubip.Get()
	if err != nil {
		cblog.Error(err)
		hostName, err := os.Hostname()
		if err != nil {
			cblog.Error(err)
		}
		return hostName
	}

	return ip.String()
}

func RunServer() {

	//======================================= setup routes
	routes := []route{
		//----------root
		{"GET", "", aw.SpiderInfo},
		{"GET", "/", aw.SpiderInfo},

		//----------CloudOS
		{"GET", "/cloudos", listCloudOS},

		//----------CloudDriverInfo
		{"POST", "/driver", registerCloudDriver},
		{"GET", "/driver", listCloudDriver},
		{"GET", "/driver/:DriverName", getCloudDriver},
		{"DELETE", "/driver/:DriverName", unRegisterCloudDriver},

		//----------CredentialInfo
		{"POST", "/credential", registerCredential},
		{"GET", "/credential", listCredential},
		{"GET", "/credential/:CredentialName", getCredential},
		{"DELETE", "/credential/:CredentialName", unRegisterCredential},

		//----------RegionInfo
		{"POST", "/region", registerRegion},
		{"GET", "/region", listRegion},
		{"GET", "/region/:RegionName", getRegion},
		{"DELETE", "/region/:RegionName", unRegisterRegion},

		//----------ConnectionConfigInfo
		{"POST", "/connectionconfig", createConnectionConfig},
		{"GET", "/connectionconfig", listConnectionConfig},
		{"GET", "/connectionconfig/:ConfigName", getConnectionConfig},
		{"DELETE", "/connectionconfig/:ConfigName", deleteConnectionConfig},

		//-------------------------------------------------------------------//

		//----------Image Handler
		{"POST", "/vmimage", createImage},
		{"GET", "/vmimage", listImage},
		{"GET", "/vmimage/:Name", getImage},
		{"DELETE", "/vmimage/:Name", deleteImage},

		//----------VMSpec Handler
		{"GET", "/vmspec", listVMSpec},
		{"GET", "/vmspec/:Name", getVMSpec},
		{"GET", "/vmorgspec", listOrgVMSpec},
		{"GET", "/vmorgspec/:Name", getOrgVMSpec},

		//----------VPC Handler
		{"POST", "/vpc", createVPC},
		{"GET", "/vpc", listVPC},
		{"GET", "/vpc/:Name", getVPC},
		{"DELETE", "/vpc/:Name", deleteVPC},
		//-- for management
		{"GET", "/allvpc", listAllVPC},
		{"DELETE", "/cspvpc/:Id", deleteCSPVPC},

		//----------SecurityGroup Handler
		{"POST", "/securitygroup", createSecurity},
		{"GET", "/securitygroup", listSecurity},
		{"GET", "/securitygroup/:Name", getSecurity},
		{"DELETE", "/securitygroup/:Name", deleteSecurity},
		//-- for management
		{"GET", "/allsecuritygroup", listAllSecurity},
		{"DELETE", "/cspsecuritygroup/:Id", deleteCSPSecurity},

		//----------KeyPair Handler
		{"POST", "/keypair", createKey},
		{"GET", "/keypair", listKey},
		{"GET", "/keypair/:Name", getKey},
		{"DELETE", "/keypair/:Name", deleteKey},
		//-- for management
		{"GET", "/allkeypair", listAllKey},
		{"DELETE", "/cspkeypair/:Id", deleteCSPKey},
		/*
			//----------VNic Handler
			{"POST", "/vnic", createVNic},
			{"GET", "/vnic", listVNic},
			{"GET", "/vnic/:VNicId", getVNic},
			{"DELETE", "/vnic/:VNicId", deleteVNic},

			//----------PublicIP Handler
			{"POST", "/publicip", createPublicIP},
			{"GET", "/publicip", listPublicIP},
			{"GET", "/publicip/:PublicIPId", getPublicIP},
			{"DELETE", "/publicip/:PublicIPId", deletePublicIP},
		*/
		//----------VM Handler
		{"POST", "/vm", startVM},
		{"GET", "/vm", listVM},
		{"GET", "/vm/:Name", getVM},
		{"DELETE", "/vm/:Name", terminateVM},
		//-- for management
		{"GET", "/allvm", listAllVM},
		{"DELETE", "/cspvm/:Id", terminateCSPVM},

		{"GET", "/vmstatus", listVMStatus},
		{"GET", "/vmstatus/:Name", getVMStatus},

		{"GET", "/controlvm/:Name", controlVM}, // suspend, resume, reboot

		//-------------------------------------------------------------------//
		//----------SSH RUN
		{"POST", "/sshrun", sshRun},

		//----------AdminWeb Handler
		{"GET", "/adminweb", aw.Frame},
		{"GET", "/adminweb/top", aw.Top},
		{"GET", "/adminweb/driver", aw.Driver},
		{"GET", "/adminweb/credential", aw.Credential},
		{"GET", "/adminweb/region", aw.Region},
		{"GET", "/adminweb/connectionconfig", aw.Connectionconfig},
		{"GET", "/adminweb/spiderinfo", aw.SpiderInfo},

		{"GET", "/adminweb/vpc/:ConnectConfig", aw.VPC},
		{"GET", "/adminweb/vpcmgmt/:ConnectConfig", aw.VPCMgmt},
		{"GET", "/adminweb/securitygroup/:ConnectConfig", aw.SecurityGroup},
		{"GET", "/adminweb/securitygroupmgmt/:ConnectConfig", aw.SecurityGroupMgmt},
		{"GET", "/adminweb/keypair/:ConnectConfig", aw.KeyPair},
		{"GET", "/adminweb/keypairmgmt/:ConnectConfig", aw.KeyPairMgmt},
		{"GET", "/adminweb/vm/:ConnectConfig", aw.VM},
		{"GET", "/adminweb/vmmgmt/:ConnectConfig", aw.VMMgmt},

		{"GET", "/adminweb/vmimage/:ConnectConfig", aw.VMImage},		
		{"GET", "/adminweb/vmspec/:ConnectConfig", aw.VMSpec},
	}
	//======================================= setup routes

	// rest's service port, now fixed.
	cr.ServicePort = ":1024"

	// Run API Server
	ApiServer(routes)

}

//================ REST API Server: setup & start
func ApiServer(routes []route) {
	e := echo.New()

	// Middleware
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	for _, route := range routes {
		// /driver => /spider/driver
		route.path = "/spider" + route.path
		switch route.method {
		case "POST":
			e.POST(route.path, route.function)
		case "GET":
			e.GET(route.path, route.function)
		case "PUT":
			e.PUT(route.path, route.function)
		case "DELETE":
			e.DELETE(route.path, route.function)

		}
	}

	// for spider logo
	cbspiderRoot := os.Getenv("CBSPIDER_ROOT")
	e.File("/spider/adminweb/images/logo.png", cbspiderRoot + "/api-runtime/rest-runtime/admin-web/images/cb-spider-circle-logo.png")

	e.HideBanner = true
	e.HidePort = true

	spiderBanner()

	e.Logger.Fatal(e.Start(cr.ServicePort))
}

//================ API Info
func apiInfo(c echo.Context) error {
        cblog.Info("call apiInfo()")

	apiInfo :=  "api info"
	return c.String(http.StatusOK, apiInfo)
}

func spiderBanner(){
	fmt.Println("\n  <CB-Spider> Multi-Cloud Infrastructure Federation Framework")

	// AdminWeb 
        adminWebURL := "http://" + cr.HostIPorName + cr.ServicePort + "/spider/adminweb"
        fmt.Printf("     - AdminWeb: %s\n", adminWebURL)

	// REST API EndPoint 
        restEndPoint := "http://" + cr.HostIPorName + cr.ServicePort + "/spider"
        fmt.Printf("     - REST API: %s\n", restEndPoint)
}
