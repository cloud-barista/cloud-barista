package cmd

import (
	"github.com/cloud-barista/cb-mcks/src/grpc-api/cbadm/app"
	"github.com/spf13/cobra"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewHealthyCmd - MCKS 상태를 수행하는 Cobra Command 생성
func NewHealthyCmd(o *app.Options) *cobra.Command {

	healthyCmd := &cobra.Command{
		Use:   "healthy",
		Short: "Healthy command for checking mcks",
		Long:  "This is a healthy command for checking mcks",
		Run: func(cmd *cobra.Command, args []string) {
			SetupAndRun(cmd, o)
		},
	}

	return healthyCmd
}
