package cbmon

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/request"
	"github.com/cloud-barista/cb-dragonfly/pkg/cmd/cbmon/cmd/get"
	"github.com/cloud-barista/cb-dragonfly/pkg/cmd/cbmon/cmd/reset"
	"github.com/cloud-barista/cb-dragonfly/pkg/cmd/cbmon/cmd/set"
	"github.com/cloud-barista/cb-dragonfly/pkg/cmd/cbmon/cmd/version"
)

// GetCLIRoot returns root command for CB-MON
func GetCLIRoot() *cobra.Command {
	// initialize cbmon command line tool
	root := &cobra.Command{
		Use:   "cbmon",
		Short: "CB-MON Command Line Interface for Cloud-Barista CB-Dragonfly framework",
	}

	// add command for cli
	root.AddCommand(
		version.NewCmd(),
		get.NewCmd(),
		set.NewCmd(),
		reset.NewCmd(),
	)

	// initialize grpc client
	monApi := request.InitMonitoringAPI()
	err := monApi.Open()
	if err != nil {
		logrus.Errorf("failed to initialize grpc client, %s", err.Error())
	}
	return root
}
