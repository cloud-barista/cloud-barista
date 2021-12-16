package cmd

import (
	"github.com/spf13/cobra"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewHealthyCmd - MCKS 상태를 수행하는 Cobra Command 생성
func NewHealthyCmd() *cobra.Command {

	healthyCmd := &cobra.Command{
		Use:   "healthy",
		Short: "This is a healthy command for checking mcks",
		Long:  "This is a healthy command for checking mcks",
		Run: func(cmd *cobra.Command, args []string) {
			SetupAndRun(cmd, args)
		},
	}

	return healthyCmd
}
