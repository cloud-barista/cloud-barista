package main

import (
	"github.com/cloud-barista/cb-milkyway/src/apiserver"
	"os"
)

func main() {

	apiserver.SPIDER_URL = os.Getenv("SPIDER_URL")

	// Run API Server
	apiserver.ApiServer()

}