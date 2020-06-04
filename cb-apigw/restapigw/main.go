package main

import (
	"log"

	"github.com/cloud-barista/cb-apigw/restapigw/cmd"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// main - Entrypoint
func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		log.Println("cb-restapigw terminated with error: ", err.Error())
	}
}

// ===== [ Public Functions ] =====
