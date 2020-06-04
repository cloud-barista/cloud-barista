package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"

	"github.com/spf13/cobra"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// runFunc - 설정파일의 문법 검증 및 구성 값 출력 및 어플리케이션 구동
func runFunc(ctx context.Context, cmd *cobra.Command, args []string) {
	var (
		sConf config.ServiceConfig
		err   error
	)

	if sConf, err = checkAndLoad(cmd, args); err != nil {
		fmt.Printf("[RUN - ERROR] %s \n", err)
		os.Exit(1)
		return
	}

	// launching the setup process
	SetupAndRun(ctx, sConf)
}

// ===== [ Public Functions ] =====

// NewRunCmd - Creates a new run command
func NewRunCmd(ctx context.Context) *cobra.Command {
	runCmd := cobra.Command{
		Use:   "run",
		Short: "Run the " + core.AppName + " server.",
		Long:  "Run the " + core.AppName + " server.",
		Run: func(cmd *cobra.Command, args []string) {
			runFunc(ctx, cmd, args)
		},
		Example: core.AppName + " run --debug --config config.yaml",
	}
	runCmd.PersistentFlags().IntVarP(&port, "port", "p", 0, "Listening port for the http service")

	return &runCmd
}
