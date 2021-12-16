// Package cmd - 어플리케이션 실행을 위한 Cobra 기반의 CLI Commands 기능 제공
package cmd

import (
	"github.com/cloud-barista/cb-mcks/src/grpc-api/config"
	"github.com/spf13/cobra"
)

// ===== [ Constants and Variables ] =====

const (
	// CLIVersion - cbadm cli 버전
	CLIVersion = "1.0"
)

var (
	configFile string
	inData     string
	inFile     string
	inType     string
	outType    string

	nameSpaceID string
	clusterName string
	nodeName    string

	parser config.Parser
)

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewRootCmd - 어플리케이션 진입점으로 사용할 Root Cobra Command 생성
func NewRootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "cbadm",
		Short: "cbadm is a lightweight grpc cli tool",
		Long:  "This is a lightweight grpc cli tool for Cloud-Barista",
	}

	// 옵션 플래그 설정
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "./grpc_conf.yaml", "config file")
	rootCmd.PersistentFlags().StringVarP(&inType, "input", "i", "yaml", "input format (json/yaml)")
	rootCmd.PersistentFlags().StringVarP(&outType, "output", "o", "yaml", "output format (json/yaml)")

	// Viper 를 사용하는 설정 파서 생성
	parser = config.MakeParser()

	//  Adds the commands for application.
	rootCmd.AddCommand(NewVersionCmd())

	rootCmd.AddCommand(NewHealthyCmd())
	rootCmd.AddCommand(NewClusterCmd())
	rootCmd.AddCommand(NewNodeCmd())

	return rootCmd
}
