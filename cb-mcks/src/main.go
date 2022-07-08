package main

import (
	"sync"

	"github.com/cloud-barista/cb-mcks/src/core/app"
	grpcserver "github.com/cloud-barista/cb-mcks/src/grpc-api/server"
	restapi "github.com/cloud-barista/cb-mcks/src/rest-api"
)

// @title CB-MCKS REST API
// @version latest
// @description CB-MCKS REST API

// @contact.name API Support
// @contact.url http://cloud-barista.github.io
// @contact.email contact-to-cloud-barista@googlegroups.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1470
// @BasePath /mcks

// @securityDefinitions.basic BasicAuth
func main() {

	app.Setup()

	wg := new(sync.WaitGroup)

	wg.Add(2)

	go func() {
		restapi.Server()
		wg.Done()
	}()

	go func() {
		grpcserver.RunServer()
		wg.Done()
	}()

	wg.Wait()

}
