package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/errors"

	"github.com/spf13/cobra"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// cheeckAndLoad - 지정된 Command 와 args 를 기준으로 시스템 설정파일 검증 및 로드
func checkAndLoad(cmd *cobra.Command, args []string) (*config.ServiceConfig, error) {
	var (
		sConf *config.ServiceConfig
		err   error
	)

	if configFile == "" {
		cmd.Println("[CHECK - ERROR] Please, provide the path to your configuration file")
		return sConf, errors.New("configuration file are not specified")
	}

	cmd.Printf("[CHECK] Parsing configuration file: %s\n", configFile)
	if sConf, err = parser.Parse(configFile); nil != err {
		return sConf, err
	}

	// Command line 에 지정된 '-d' 옵션 우선 적용
	sConf.Debug = sConf.Debug || debug

	// Command line 에 지정된 '-p' 옵션 우선 적용
	if port != 0 {
		sConf.Port = port
	}

	// 서비스 설정 출력 (JSON)
	cmd.Println("[CHECK - SYSTEM CONFIGURATION] \n", core.ToJSON(sConf))

	// 서비스 설정 검증
	if err = sConf.Validate(); nil != err {
		return nil, err
	}

	return sConf, err
}

// checkFunc - 지정된 args 에서 설정과 관련된 정보를 로드/검증/출력 처리
func checkFunc(cmd *cobra.Command, args []string) {
	var (
		err error
	)

	if _, err = checkAndLoad(cmd, args); nil != err {
		fmt.Printf("[CHECK - ERROR] %s \n", err)
		os.Exit(1)
		return
	}

	cmd.Println("[CHECK] Syntax OK!")
}

// ===== [ Public Functions ] =====

// NewCheckCmd - 설정 검증 기능을 수행하는 Cobra Command 생성
func NewCheckCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "check",
		Short:   "Validates that the configuration file is valid",
		Long:    "Validates that the active configuration file has  a valid syntax to run the service. \nChange the configuration file by using the --config flag (default $PWD/conf/cb-restapigw.yaml)",
		Run:     checkFunc,
		Aliases: []string{"validate"},
		Example: core.AppName + " check --config config.yaml",
	}
}
