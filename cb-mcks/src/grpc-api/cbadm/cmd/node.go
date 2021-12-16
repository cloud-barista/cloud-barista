package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cloud-barista/cb-mcks/src/grpc-api/logger"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewNodeCmd - Node 관리 기능을 수행하는 Cobra Command 생성
func NewNodeCmd() *cobra.Command {

	nodeCmd := &cobra.Command{
		Use:   "node",
		Short: "This is a manageable command for node",
		Long:  "This is a manageable command for node",
	}

	//  Adds the commands for application.
	nodeCmd.AddCommand(NewNodeAddCmd())
	nodeCmd.AddCommand(NewNodeListCmd())
	nodeCmd.AddCommand(NewNodeGetCmd())
	nodeCmd.AddCommand(NewNodeRemoveCmd())

	return nodeCmd
}

// NewNodeAddCmd - Node 생성 기능을 수행하는 Cobra Command 생성
func NewNodeAddCmd() *cobra.Command {

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "This is add command for node",
		Long:  "This is add command for node",
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

// NewNodeListCmd - Node 목록 기능을 수행하는 Cobra Command 생성
func NewNodeListCmd() *cobra.Command {

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "This is list command for node",
		Long:  "This is list command for node",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if nameSpaceID == "" {
				logger.Error("failed to validate --ns parameter")
				return
			}
			if clusterName == "" {
				logger.Error("failed to validate --cluster parameter")
				return
			}

			logger.Debug("--ns parameter value : ", nameSpaceID)
			logger.Debug("--cluster parameter value : ", clusterName)

			SetupAndRun(cmd, args)
		},
	}

	listCmd.PersistentFlags().StringVarP(&nameSpaceID, "ns", "", "", "namespace id")
	listCmd.PersistentFlags().StringVarP(&clusterName, "cluster", "", "", "cluster name")

	return listCmd
}

// NewNodeGetCmd - Node 조회 기능을 수행하는 Cobra Command 생성
func NewNodeGetCmd() *cobra.Command {

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "This is get command for node",
		Long:  "This is get command for node",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if nameSpaceID == "" {
				logger.Error("failed to validate --ns parameter")
				return
			}
			if clusterName == "" {
				logger.Error("failed to validate --cluster parameter")
				return
			}
			if nodeName == "" {
				logger.Error("failed to validate --node parameter")
				return
			}
			logger.Debug("--ns parameter value : ", nameSpaceID)
			logger.Debug("--cluster parameter value : ", clusterName)
			logger.Debug("--node parameter value : ", nodeName)

			SetupAndRun(cmd, args)
		},
	}

	getCmd.PersistentFlags().StringVarP(&nameSpaceID, "ns", "", "", "namespace id")
	getCmd.PersistentFlags().StringVarP(&clusterName, "cluster", "", "", "cluster name")
	getCmd.PersistentFlags().StringVarP(&nodeName, "node", "", "", "node name")

	return getCmd
}

// NewNodeRemoveCmd - Node 삭제 기능을 수행하는 Cobra Command 생성
func NewNodeRemoveCmd() *cobra.Command {

	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "This is remove command for node",
		Long:  "This is remove command for node",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if nameSpaceID == "" {
				logger.Error("failed to validate --ns parameter")
				return
			}
			if clusterName == "" {
				logger.Error("failed to validate --cluster parameter")
				return
			}
			if nodeName == "" {
				logger.Error("failed to validate --node parameter")
				return
			}
			logger.Debug("--ns parameter value : ", nameSpaceID)
			logger.Debug("--cluster parameter value : ", clusterName)
			logger.Debug("--node parameter value : ", nodeName)

			SetupAndRun(cmd, args)
		},
	}

	removeCmd.PersistentFlags().StringVarP(&nameSpaceID, "ns", "", "", "namespace id")
	removeCmd.PersistentFlags().StringVarP(&clusterName, "cluster", "", "", "cluster name")
	removeCmd.PersistentFlags().StringVarP(&nodeName, "node", "", "", "node name")

	return removeCmd
}
