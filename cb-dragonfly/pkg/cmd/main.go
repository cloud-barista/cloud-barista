package main

import (
	"os"

	"github.com/cloud-barista/cb-dragonfly/pkg/cmd/cbmon"
)

func main() {
	root := cbmon.GetCLIRoot()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
