// gRPC Runtime of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by CB-Spider Team, 2020.09.

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cloud-barista/cb-spider/api-runtime/grpc-runtime/logger"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewKeyPairCmd - KeyPair 관리 기능을 수행하는 Cobra Command 생성
func NewKeyPairCmd() *cobra.Command {

	keyPairCmd := &cobra.Command{
		Use:   "keypair",
		Short: "This is a manageable command for keypair",
		Long:  "This is a manageable command for keypair",
	}

	//  Adds the commands for application.
	keyPairCmd.AddCommand(NewKeyPairCreateCmd())
	keyPairCmd.AddCommand(NewKeyPairListCmd())
	keyPairCmd.AddCommand(NewKeyPairGetCmd())
	keyPairCmd.AddCommand(NewKeyPairDeleteCmd())
	keyPairCmd.AddCommand(NewKeyPairListAllCmd())
	keyPairCmd.AddCommand(NewKeyPairDeleteCSPCmd())

	return keyPairCmd
}

// NewKeyPairCreateCmd - KeyPair 생성 기능을 수행하는 Cobra Command 생성
func NewKeyPairCreateCmd() *cobra.Command {

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "This is create command for keypair",
		Long:  "This is create command for keypair",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			readInDataFromFile()
			if inData == "" {
				logger.Error("failed to validate --indata parameter")
				return
			}
			logger.Debug("--indata parameter value : \n", inData)
			logger.Debug("--infile parameter value : ", inFile)

			SetupAndRun(cmd, args)
		},
	}

	createCmd.PersistentFlags().StringVarP(&inData, "indata", "d", "", "input string data")
	createCmd.PersistentFlags().StringVarP(&inFile, "infile", "f", "", "input file path")

	return createCmd
}

// NewKeyPairListCmd - KeyPair 목록 기능을 수행하는 Cobra Command 생성
func NewKeyPairListCmd() *cobra.Command {

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "This is list command for keypair",
		Long:  "This is list command for keypair",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if connectionName == "" {
				logger.Error("failed to validate --cname parameter")
				return
			}
			logger.Debug("--cname parameter value : ", connectionName)

			SetupAndRun(cmd, args)
		},
	}

	listCmd.PersistentFlags().StringVarP(&connectionName, "cname", "", "", "connection name")

	return listCmd
}

// NewKeyPairGetCmd - KeyPair 조회 기능을 수행하는 Cobra Command 생성
func NewKeyPairGetCmd() *cobra.Command {

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "This is get command for keypair",
		Long:  "This is get command for keypair",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if connectionName == "" {
				logger.Error("failed to validate --cname parameter")
				return
			}
			if keypairName == "" {
				logger.Error("failed to validate --name parameter")
				return
			}
			logger.Debug("--cname parameter value : ", connectionName)
			logger.Debug("--name parameter value : ", keypairName)

			SetupAndRun(cmd, args)
		},
	}

	getCmd.PersistentFlags().StringVarP(&connectionName, "cname", "", "", "connection name")
	getCmd.PersistentFlags().StringVarP(&keypairName, "name", "n", "", "keypair name")

	return getCmd
}

// NewKeyPairDeleteCmd - KeyPair 삭제 기능을 수행하는 Cobra Command 생성
func NewKeyPairDeleteCmd() *cobra.Command {

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "This is delete command for keypair",
		Long:  "This is delete command for keypair",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if connectionName == "" {
				logger.Error("failed to validate --cname parameter")
				return
			}
			if keypairName == "" {
				logger.Error("failed to validate --name parameter")
				return
			}
			logger.Debug("--cname parameter value : ", connectionName)
			logger.Debug("--name parameter value : ", keypairName)
			logger.Debug("--force parameter value : ", force)

			SetupAndRun(cmd, args)
		},
	}

	deleteCmd.PersistentFlags().StringVarP(&connectionName, "cname", "", "", "connection name")
	deleteCmd.PersistentFlags().StringVarP(&keypairName, "name", "n", "", "keypair name")
	deleteCmd.PersistentFlags().StringVarP(&force, "force", "", "false", "force flg (true/false)")

	return deleteCmd
}

// NewKeyPairListAllCmd - 관리 KeyPair 목록 기능을 수행하는 Cobra Command 생성
func NewKeyPairListAllCmd() *cobra.Command {

	listAllCmd := &cobra.Command{
		Use:   "listall",
		Short: "This is list all command for keypair",
		Long:  "This is list all command for keypair",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if connectionName == "" {
				logger.Error("failed to validate --cname parameter")
				return
			}
			logger.Debug("--cname parameter value : ", connectionName)

			SetupAndRun(cmd, args)
		},
	}

	listAllCmd.PersistentFlags().StringVarP(&connectionName, "cname", "", "", "connection name")

	return listAllCmd
}

// NewKeyPairDeleteCSPCmd - 관리 KeyPair 삭제 기능을 수행하는 Cobra Command 생성
func NewKeyPairDeleteCSPCmd() *cobra.Command {

	deleteCSPCmd := &cobra.Command{
		Use:   "deletecsp",
		Short: "This is delete csp command for keypair",
		Long:  "This is delete csp command for keypair",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if connectionName == "" {
				logger.Error("failed to validate --cname parameter")
				return
			}
			if cspID == "" {
				logger.Error("failed to validate --id parameter")
				return
			}
			logger.Debug("--cname parameter value : ", connectionName)
			logger.Debug("--id parameter value : ", cspID)

			SetupAndRun(cmd, args)
		},
	}

	deleteCSPCmd.PersistentFlags().StringVarP(&connectionName, "cname", "", "", "connection name")
	deleteCSPCmd.PersistentFlags().StringVarP(&cspID, "id", "", "", "csp id")

	return deleteCSPCmd
}
