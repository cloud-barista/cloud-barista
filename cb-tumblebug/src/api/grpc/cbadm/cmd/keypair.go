package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cloud-barista/cb-tumblebug/src/api/grpc/logger"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewKeypairCmd : "cbadm keypair *" (for CB-Tumblebug)
func NewKeypairCmd() *cobra.Command {

	keypairCmd := &cobra.Command{
		Use:   "keypair",
		Short: "This is a manageable command for keypair",
		Long:  "This is a manageable command for keypair",
	}

	//  Adds the commands for application.
	keypairCmd.AddCommand(NewKeypairCreateCmd())
	keypairCmd.AddCommand(NewKeypairListCmd())
	keypairCmd.AddCommand(NewKeypairListIdCmd())
	keypairCmd.AddCommand(NewKeypairGetCmd())
	keypairCmd.AddCommand(NewKeypairSaveCmd())
	keypairCmd.AddCommand(NewKeypairDeleteCmd())
	keypairCmd.AddCommand(NewKeypairDeleteAllCmd())

	return keypairCmd
}

// NewKeypairCreateCmd : "cbadm keypair create"
func NewKeypairCreateCmd() *cobra.Command {

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

// NewKeypairListCmd : "cbadm keypair list"
func NewKeypairListCmd() *cobra.Command {

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "This is list command for keypair",
		Long:  "This is list command for keypair",
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

// NewKeypairListIdCmd : "cbadm keypair list-id"
func NewKeypairListIdCmd() *cobra.Command {

	listIdCmd := &cobra.Command{
		Use:   "list-id",
		Short: "This is list-id command for keypair",
		Long:  "This is list-id command for keypair",
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

	listIdCmd.PersistentFlags().StringVarP(&nameSpaceID, "ns", "", "", "namespace id")

	return listIdCmd
}

// NewKeypairGetCmd : "cbadm keypair get"
func NewKeypairGetCmd() *cobra.Command {

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "This is get command for keypair",
		Long:  "This is get command for keypair",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if nameSpaceID == "" {
				logger.Error("failed to validate --ns parameter")
				return
			}
			if resourceID == "" {
				logger.Error("failed to validate --id parameter")
				return
			}
			logger.Debug("--ns parameter value : ", nameSpaceID)
			logger.Debug("--id parameter value : ", resourceID)

			SetupAndRun(cmd, args)
		},
	}

	getCmd.PersistentFlags().StringVarP(&nameSpaceID, "ns", "", "", "namespace id")
	getCmd.PersistentFlags().StringVarP(&resourceID, "id", "", "", "keypair id")

	return getCmd
}

// NewKeypairSaveCmd : "cbadm keypair save"
func NewKeypairSaveCmd() *cobra.Command {

	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "This is save command for keypair",
		Long:  "This is save command for keypair",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if nameSpaceID == "" {
				logger.Error("failed to validate --ns parameter")
				return
			}
			if resourceID == "" {
				logger.Error("failed to validate --id parameter")
				return
			}
			if sshSaveFileName == "" {
				logger.Error("failed to validate --fn parameter")
				return
			}
			logger.Debug("--ns parameter value : ", nameSpaceID)
			logger.Debug("--id parameter value : ", resourceID)
			logger.Debug("--fn parameter value : ", sshSaveFileName)

			SetupAndRun(cmd, args)
		},
	}

	saveCmd.PersistentFlags().StringVarP(&nameSpaceID, "ns", "", "", "namespace id")
	saveCmd.PersistentFlags().StringVarP(&resourceID, "id", "", "", "keypair id")
	saveCmd.PersistentFlags().StringVarP(&sshSaveFileName, "fn", "", "", "ssh key save file name")

	return saveCmd
}

// NewKeypairDeleteCmd : "cbadm keypair delete"
func NewKeypairDeleteCmd() *cobra.Command {

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "This is delete command for keypair",
		Long:  "This is delete command for keypair",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if nameSpaceID == "" {
				logger.Error("failed to validate --ns parameter")
				return
			}
			if resourceID == "" {
				logger.Error("failed to validate --id parameter")
				return
			}
			if force == "" {
				logger.Error("failed to validate --force parameter")
				return
			}
			logger.Debug("--ns parameter value : ", nameSpaceID)
			logger.Debug("--id parameter value : ", resourceID)
			logger.Debug("--force parameter value : ", force)

			SetupAndRun(cmd, args)
		},
	}

	deleteCmd.PersistentFlags().StringVarP(&nameSpaceID, "ns", "", "", "namespace id")
	deleteCmd.PersistentFlags().StringVarP(&resourceID, "id", "", "", "keypair id")
	deleteCmd.PersistentFlags().StringVarP(&force, "force", "", "false", "force flag")

	return deleteCmd
}

// NewKeypairDeleteAllCmd : "cbadm keypair delete-all"
func NewKeypairDeleteAllCmd() *cobra.Command {

	deleteAllCmd := &cobra.Command{
		Use:   "delete-all",
		Short: "This is delete-all command for keypair",
		Long:  "This is delete-all command for keypair",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewLogger()
			if nameSpaceID == "" {
				logger.Error("failed to validate --ns parameter")
				return
			}
			if force == "" {
				logger.Error("failed to validate --force parameter")
				return
			}
			logger.Debug("--ns parameter value : ", nameSpaceID)
			logger.Debug("--force parameter value : ", force)

			SetupAndRun(cmd, args)
		},
	}

	deleteAllCmd.PersistentFlags().StringVarP(&nameSpaceID, "ns", "", "", "namespace id")
	deleteAllCmd.PersistentFlags().StringVarP(&force, "force", "", "false", "force flag")

	return deleteAllCmd
}
