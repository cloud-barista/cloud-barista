// Proof of Concepts for the Cloud-Barista Multi-Cloud Project.
//      * Cloud-Barista: https://github.com/cloud-barista

package apiserver

import (
	"github.com/cloud-barista/cb-milkyway/src/common"

	//"os"

	"fmt"

	// REST API (echo)
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	// CB-Store
)

/*
// CB-Store
var cblog *logrus.Logger
var store icbs.Store

func init() {
	cblog = config.Cblogger
	store = cbstore.GetStore()
}

type KeyValue struct {
	Key   string
	Value string
}
*/

//var masterConfigInfos confighandler.MASTERCONFIGTYPE

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

const (
	Version = " Version: Cappuccino"
	website = " Repository: https://github.com/cloud-barista/cb-milkyway"
	banner  = `

 ██████╗██████╗       ███╗   ███╗██╗██╗     ██╗  ██╗██╗   ██╗██╗    ██╗ █████╗ ██╗   ██╗
██╔════╝██╔══██╗      ████╗ ████║██║██║     ██║ ██╔╝╚██╗ ██╔╝██║    ██║██╔══██╗╚██╗ ██╔╝
██║     ██████╔╝█████╗██╔████╔██║██║██║     █████╔╝  ╚████╔╝ ██║ █╗ ██║███████║ ╚████╔╝ 
██║     ██╔══██╗╚════╝██║╚██╔╝██║██║██║     ██╔═██╗   ╚██╔╝  ██║███╗██║██╔══██║  ╚██╔╝  
╚██████╗██████╔╝      ██║ ╚═╝ ██║██║███████╗██║  ██╗   ██║   ╚███╔███╔╝██║  ██║   ██║   
 ╚═════╝╚═════╝       ╚═╝     ╚═╝╚═╝╚══════╝╚═╝  ╚═╝   ╚═╝    ╚══╝╚══╝ ╚═╝  ╚═╝   ╚═╝                    

 Benchmark Agent for CB-Tumblebug
 ________________________________________________`
)

// Main Body

func ApiServer() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! This is cb-tumblebug Agent")
	})
	e.HideBanner = true
	//e.colorer.Printf(banner, e.colorer.Red("v"+Version), e.colorer.Blue(website))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Printf(banner)
	fmt.Println("")
	fmt.Printf(ErrorColor, Version)
	fmt.Println("")
	fmt.Printf(InfoColor, website)
	fmt.Println("")
	fmt.Println("")

	// Route
	g := e.Group("/milkyway", common.ApiValidation())

	//g.POST("", common.RestPostNs)
	//g.GET("/:nsId", common.RestGetNs)

	g.GET("/install", common.RestGetInstall)
	g.GET("/init", common.RestGetInit)
	g.GET("/clean", common.RestGetClean)

	g.GET("/cpus", common.RestGetCPUS)
	g.GET("/cpum", common.RestGetCPUM)
	g.GET("/memR", common.RestGetMEMR)
	g.GET("/memW", common.RestGetMEMW)
	g.GET("/fioR", common.RestGetFIOR)
	g.GET("/fioW", common.RestGetFIOW)
	g.GET("/dbR", common.RestGetDBR)
	g.GET("/dbW", common.RestGetDBW)

	g.GET("/rtt", common.RestGetRTT)
	g.GET("/mrtt", common.RestGetMultiRTT)
	
	//g.PUT("/:nsId", common.RestPutNs)
	//g.DELETE("/:nsId", common.RestDelNs)
	//g.DELETE("", common.RestDelAllNs)

	/*
	g.POST("", common.RestPostNs)
	g.GET("/:nsId", common.RestGetNs)
	g.GET("", common.RestGetAllNs)
	g.PUT("/:nsId", common.RestPutNs)
	g.DELETE("/:nsId", common.RestDelNs)
	g.DELETE("", common.RestDelAllNs)
	*/

	e.Logger.Fatal(e.Start(":1324"))

}

var SPIDER_URL string
