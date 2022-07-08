// Package cmd - 어플리케이션 실행을 위한 Cobra 기반의 CLI Commands 기능 제공
package cmd

import (
	"os"

	"github.com/cloud-barista/cb-mcks/src/grpc-api/cbadm/app"
	"github.com/spf13/cobra"
)

// ===== [ Constants and Variables ] =====

const (
	// CLIVersion - cbadm cli 버전
	CLIVersion = "1.0"
)

var (
	clusterName string
)

type CbadmOptions struct {
	app.Options
}

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewRootCmd - 어플리케이션 진입점으로 사용할 Root Cobra Command 생성
func NewRootCmd() *cobra.Command {

	o := CbadmOptions{
		Options: app.Options{
			OutStream: os.Stdout,
		},
	}

	rootCmd := &cobra.Command{
		Use:   "cbadm",
		Short: "cbadm is a lightweight grpc cli tool",
		Long:  "This is a lightweight grpc cli tool for Cloud-Barista",
	}

	// 옵션 플래그 설정
	rootCmd.PersistentFlags().StringVar(&o.Name, "name", "", "name")
	rootCmd.PersistentFlags().StringVarP(&o.ConfigFile, "config", "c", "", "configuration file path")
	rootCmd.PersistentFlags().StringVarP(&o.Namespace, "namespace", "n", "", "cloud-baristar namespace")
	rootCmd.PersistentFlags().StringVarP(&o.Filename, "file", "f", "", "filepath")
	rootCmd.PersistentFlags().StringVarP(&o.Data, "data", "d", "", "input string data")
	rootCmd.PersistentFlags().StringVarP(&o.Output, "output", "o", "yaml", "output format (json/yaml)")

	if err := app.OnConfigInitialize(o.ConfigFile); err != nil {
		o.PrintlnError(err)
		os.Exit(1)
	}

	//  Adds the commands for application.
	rootCmd.AddCommand(NewCommandConfig(&o.Options))
	rootCmd.AddCommand(NewVersionCmd())

	rootCmd.AddCommand(NewHealthyCmd(&o.Options))
	rootCmd.AddCommand(NewGetCmd(&o.Options))
	rootCmd.AddCommand(NewCreateCmd(&o.Options))
	rootCmd.AddCommand(NewDeleteCmd(&o.Options))

	return rootCmd
}
