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

// NewVPCCmd - VPC 관리 기능을 수행하는 Cobra Command 생성
func NewVPCCmd() *cobra.Command {

	vpcCmd := &cobra.Command{
		Use:   "vpc",
		Short: "This is a manageable command for vpc",
		Long:  "This is a manageable command for vpc",
	}

	//  Adds the commands for application.
	vpcCmd.AddCommand(NewVPCCreateCmd())
	vpcCmd.AddCommand(NewVPCListCmd())
	vpcCmd.AddCommand(NewVPCGetCmd())
	vpcCmd.AddCommand(NewVPCDeleteCmd())
	vpcCmd.AddCommand(NewVPCListAllCmd())
	vpcCmd.AddCommand(NewVPCDeleteCSPCmd())
	vpcCmd.AddCommand(NewSubnetAddCmd())
	vpcCmd.AddCommand(NewSubnetRemoveCmd())
	vpcCmd.AddCommand(NewSubnetRemoveCSPCmd())

	return vpcCmd
}

// NewVPCCreateCmd - VPC 생성 기능을 수행하는 Cobra Command 생성
func NewVPCCreateCmd() *cobra.Command {

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "This is create command for vpc",
		Long:  "This is create command for vpc",
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

// NewVPCListCmd - VPC 목록 기능을 수행하는 Cobra Command 생성
func NewVPCListCmd() *cobra.Command {

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "This is list command for vpc",
		Long:  "This is list command for vpc",
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

// NewVPCGetCmd - VPC 조회 기능을 수행하는 Cobra Command 생성
func NewVPCGetCmd() *cobra.Command {

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "This is get command for vpc",
		Long:  "This is get command for vpc",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if connectionName == "" {
				logger.Error("failed to validate --cname parameter")
				return
			}
			if vpcName == "" {
				logger.Error("failed to validate --name parameter")
				return
			}
			logger.Debug("--cname parameter value : ", connectionName)
			logger.Debug("--name parameter value : ", vpcName)

			SetupAndRun(cmd, args)
		},
	}

	getCmd.PersistentFlags().StringVarP(&connectionName, "cname", "", "", "connection name")
	getCmd.PersistentFlags().StringVarP(&vpcName, "name", "n", "", "vpc name")

	return getCmd
}

// NewVPCDeleteCmd - VPC 삭제 기능을 수행하는 Cobra Command 생성
func NewVPCDeleteCmd() *cobra.Command {

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "This is delete command for vpc",
		Long:  "This is delete command for vpc",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if connectionName == "" {
				logger.Error("failed to validate --cname parameter")
				return
			}
			if vpcName == "" {
				logger.Error("failed to validate --name parameter")
				return
			}
			logger.Debug("--cname parameter value : ", connectionName)
			logger.Debug("--name parameter value : ", vpcName)
			logger.Debug("--force parameter value : ", force)

			SetupAndRun(cmd, args)
		},
	}

	deleteCmd.PersistentFlags().StringVarP(&connectionName, "cname", "", "", "connection name")
	deleteCmd.PersistentFlags().StringVarP(&vpcName, "name", "n", "", "vpc name")
	deleteCmd.PersistentFlags().StringVarP(&force, "force", "", "false", "force flg (true/false)")

	return deleteCmd
}

// NewVPCListAllCmd - 관리 VPC 목록 기능을 수행하는 Cobra Command 생성
func NewVPCListAllCmd() *cobra.Command {

	listAllCmd := &cobra.Command{
		Use:   "listall",
		Short: "This is list all command for vpc",
		Long:  "This is list all command for vpc",
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

// NewVPCDeleteCSPCmd - 관리 VPC 삭제 기능을 수행하는 Cobra Command 생성
func NewVPCDeleteCSPCmd() *cobra.Command {

	deleteCSPCmd := &cobra.Command{
		Use:   "deletecsp",
		Short: "This is delete csp command for vpc",
		Long:  "This is delete csp command for vpc",
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

// NewSubnetAddCmd - Subnet 추가 기능을 수행하는 Cobra Command 생성
func NewSubnetAddCmd() *cobra.Command {

	addCmd := &cobra.Command{
		Use:   "add-subnet",
		Short: "This is add command for vpc subnet",
		Long:  "This is add command for vpc subnet",
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

	addCmd.PersistentFlags().StringVarP(&inData, "indata", "d", "", "input string data")
	addCmd.PersistentFlags().StringVarP(&inFile, "infile", "f", "", "input file path")

	return addCmd
}

// NewSubnetRemoveCmd - Subnet 삭제 기능을 수행하는 Cobra Command 생성
func NewSubnetRemoveCmd() *cobra.Command {

	removeCmd := &cobra.Command{
		Use:   "remove-subnet",
		Short: "This is remove command for vpc subnet",
		Long:  "This is remove command for vpc subnet",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if connectionName == "" {
				logger.Error("failed to validate --cname parameter")
				return
			}
			if vpcName == "" {
				logger.Error("failed to validate --vname parameter")
				return
			}
			if subnetName == "" {
				logger.Error("failed to validate --sname parameter")
				return
			}
			logger.Debug("--cname parameter value : ", connectionName)
			logger.Debug("--vname parameter value : ", vpcName)
			logger.Debug("--sname parameter value : ", subnetName)
			logger.Debug("--force parameter value : ", force)

			SetupAndRun(cmd, args)
		},
	}

	removeCmd.PersistentFlags().StringVarP(&connectionName, "cname", "", "", "connection name")
	removeCmd.PersistentFlags().StringVarP(&vpcName, "vname", "", "", "vpc name")
	removeCmd.PersistentFlags().StringVarP(&subnetName, "sname", "", "", "subnet name")
	removeCmd.PersistentFlags().StringVarP(&force, "force", "", "false", "force flg (true/false)")

	return removeCmd
}

// NewSubnetRemoveCSPCmd - CSP Subnet 삭제 기능을 수행하는 Cobra Command 생성
func NewSubnetRemoveCSPCmd() *cobra.Command {

	removeCSPCmd := &cobra.Command{
		Use:   "removecsp-subnet",
		Short: "This is remove csp command for vpc subnet",
		Long:  "This is remove csp command for vpc subnet",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if connectionName == "" {
				logger.Error("failed to validate --cname parameter")
				return
			}
			if vpcName == "" {
				logger.Error("failed to validate --vname parameter")
				return
			}
			if cspID == "" {
				logger.Error("failed to validate --id parameter")
				return
			}
			logger.Debug("--cname parameter value : ", connectionName)
			logger.Debug("--vname parameter value : ", vpcName)
			logger.Debug("--id parameter value : ", cspID)

			SetupAndRun(cmd, args)
		},
	}

	removeCSPCmd.PersistentFlags().StringVarP(&connectionName, "cname", "", "", "connection name")
	removeCSPCmd.PersistentFlags().StringVarP(&vpcName, "vname", "", "", "vpc name")
	removeCSPCmd.PersistentFlags().StringVarP(&cspID, "id", "", "", "csp id")

	return removeCSPCmd
}
