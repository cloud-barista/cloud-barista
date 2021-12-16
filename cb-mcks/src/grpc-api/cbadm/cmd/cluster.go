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

// NewClusterCmd - Cluster 관리 기능을 수행하는 Cobra Command 생성
func NewClusterCmd() *cobra.Command {

	clusterCmd := &cobra.Command{
		Use:   "cluster",
		Short: "This is a manageable command for cluster",
		Long:  "This is a manageable command for cluster",
	}

	//  Adds the commands for application.
	clusterCmd.AddCommand(NewClusterCreateCmd())
	clusterCmd.AddCommand(NewClusterListCmd())
	clusterCmd.AddCommand(NewClusterGetCmd())
	clusterCmd.AddCommand(NewClusterDeleteCmd())

	return clusterCmd
}

// NewClusterCreateCmd - Cluster 생성 기능을 수행하는 Cobra Command 생성
func NewClusterCreateCmd() *cobra.Command {

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "This is create command for cluster",
		Long:  "This is create command for cluster",
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

// NewClusterListCmd - Cluster 목록 기능을 수행하는 Cobra Command 생성
func NewClusterListCmd() *cobra.Command {

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "This is list command for cluster",
		Long:  "This is list command for cluster",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if nameSpaceID == "" {
				logger.Error("failed to validate --ns parameter")
				return
			}
			logger.Debug("--ns parameter value : ", nameSpaceID)

			SetupAndRun(cmd, args)
		},
	}

	listCmd.PersistentFlags().StringVarP(&nameSpaceID, "ns", "", "", "namespace id")

	return listCmd
}

// NewClusterGetCmd - Cluster 조회 기능을 수행하는 Cobra Command 생성
func NewClusterGetCmd() *cobra.Command {

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "This is get command for cluster",
		Long:  "This is get command for cluster",
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

	getCmd.PersistentFlags().StringVarP(&nameSpaceID, "ns", "", "", "namespace id")
	getCmd.PersistentFlags().StringVarP(&clusterName, "cluster", "", "", "cluster name")

	return getCmd
}

// NewClusterDeleteCmd - Cluster 삭제 기능을 수행하는 Cobra Command 생성
func NewClusterDeleteCmd() *cobra.Command {

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "This is delete command for cluster",
		Long:  "This is delete command for cluster",
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

	deleteCmd.PersistentFlags().StringVarP(&nameSpaceID, "ns", "", "", "namespace id")
	deleteCmd.PersistentFlags().StringVarP(&clusterName, "cluster", "", "", "cluster name")

	return deleteCmd
}
