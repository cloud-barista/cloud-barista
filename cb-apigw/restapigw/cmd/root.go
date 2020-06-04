// Package cmd - 어플리케이션 실행을 위한 Cobra 기반의 CLI Commands 기능 제공
package cmd

import (
	"context"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"

	"github.com/spf13/cobra"
)

// ===== [ Constants and Variables ] =====

var (
	configFile string
	debug      bool
	port       int
	parser     config.Parser
)

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewRootCmd - 어플리케이션 진입점으로 사용할 Root Cobra Command 생성
func NewRootCmd() *cobra.Command {
	ctx := context.Background()

	cmd := &cobra.Command{
		Use:   core.AppName,
		Short: "`" + core.AppName + "` is an RESTful API Gateway",
		Long:  "This is a lightweight RESTful API Gateway and Management Platform for Cloud-Barista",
	}

	// 옵션 플래그 설정
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file (default is $PWD/conf/cb-restapigw.yaml")
	cmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable the debug mode")

	// Viper 를 사용하는 설정 파서 생성
	parser = config.MakeParser()

	//  Adds the commands for check and run application.
	cmd.AddCommand(NewCheckCmd(ctx))
	cmd.AddCommand(NewRunCmd(ctx))

	return cmd
}
