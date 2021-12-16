package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/cloud-barista/cb-mcks/src/grpc-api/logger"
	lb_api "github.com/cloud-barista/cb-mcks/src/grpc-api/request"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

func readInDataFromFile() {
	logger := logger.NewLogger()
	if inData == "" {
		if inFile != "" {
			dat, err := ioutil.ReadFile(inFile)
			if err != nil {
				logger.Error("failed to read file : ", inFile)
				return
			}
			inData = string(dat)
		}
	}
}

// ===== [ Public Functions ] =====

// SetupAndRun - MCKS GRPC CLI 구동
func SetupAndRun(cmd *cobra.Command, args []string) {
	logger := logger.NewLogger()

	var (
		result string
		err    error

		mcar *lb_api.MCARApi = nil
	)

	// panic 처리
	defer func() {
		if r := recover(); r != nil {
			logger.Error("cbadm is stopped : ", r)
		}
	}()

	if cmd.Parent().Name() == "cluster" || cmd.Parent().Name() == "node" || cmd.Name() == "healthy" {
		// LB API 설정
		mcar = lb_api.NewMCARManager()
		err = mcar.SetConfigPath(configFile)
		if err != nil {
			logger.Error("failed to set config : ", err)
			return
		}
		err = mcar.Open()
		if err != nil {
			logger.Error("mcks api open failed : ", err)
			return
		}
		defer mcar.Close()
	}

	// 입력 파라미터 처리
	if outType != "json" && outType != "yaml" {
		logger.Error("failed to validate --output parameter : ", outType)
		return
	}
	if inType != "json" && inType != "yaml" {
		logger.Error("failed to validate --input parameter : ", inType)
		return
	}

	if cmd.Parent().Name() == "cluster" || cmd.Parent().Name() == "node" || cmd.Name() == "healthy" {
		mcar.SetInType(inType)
		mcar.SetOutType(outType)
	}

	logger.Debug("--input parameter value : ", inType)
	logger.Debug("--output parameter value : ", outType)

	result = ""
	err = nil

	switch cmd.Parent().Name() {
	case "cbadm":
		switch cmd.Name() {
		case "healthy":
			result, err = mcar.Healthy()
		}
	case "cluster":
		switch cmd.Name() {
		case "create":
			result, err = mcar.CreateCluster(inData)
		case "list":
			result, err = mcar.ListClusterByParam(nameSpaceID)
		case "get":
			result, err = mcar.GetClusterByParam(nameSpaceID, clusterName)
		case "delete":
			result, err = mcar.DeleteClusterByParam(nameSpaceID, clusterName)
		}
	case "node":
		switch cmd.Name() {
		case "add":
			result, err = mcar.AddNode(inData)
		case "list":
			result, err = mcar.ListNodeByParam(nameSpaceID, clusterName)
		case "get":
			result, err = mcar.GetNodeByParam(nameSpaceID, clusterName, nodeName)
		case "remove":
			result, err = mcar.RemoveNodeByParam(nameSpaceID, clusterName, nodeName)
		}
	}

	if err != nil {
		if outType == "yaml" {
			fmt.Fprintf(cmd.OutOrStdout(), "message: %v\n", err)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "{\"message\": \"%v\"}\n", err)
		}
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", result)
	}

}
