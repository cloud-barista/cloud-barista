package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cloud-barista/cb-mcks/src/grpc-api/cbadm/app"
	"github.com/cloud-barista/cb-mcks/src/grpc-api/logger"
	lb_api "github.com/cloud-barista/cb-mcks/src/grpc-api/request"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// SetupAndRun - MCKS GRPC CLI 구동
func SetupAndRun(cmd *cobra.Command, o *app.Options) {
	logger := logger.NewLogger()

	var (
		result = ""
		err    error
	)

	// panic 처리
	defer func() {
		if r := recover(); r != nil {
			logger.Error("cbadm is stopped : ", r)
		}
	}()

	if o.Output != "json" && o.Output != "yaml" {
		logger.Error("failed to validate --output parameter : ", o.Output)
		return
	}
	mcar := lb_api.NewMCARManager()
	//cim := sp_api.NewCloudInfoManager()

	if cmd.Name() == "cluster" || cmd.Name() == "node" || cmd.Name() == "healthy" {
		// LB API 설정
		mckscli := app.Config.GetCurrentContext().Mckscli

		err := mcar.SetServerAddr(mckscli.ServerAddr)
		if err != nil {
			logger.Error("server_addr set failed", err)
		}

		timeout, _ := time.ParseDuration(mckscli.Timeout)
		err = mcar.SetTimeout(timeout)
		if err != nil {
			logger.Error("timeout set failed", err)
		}
		err = mcar.Open()
		if err != nil {
			logger.Error("mcks api open failed : ", err)
			return
		}
		defer mcar.Close()

		mcar.SetInType("json")
		mcar.SetOutType(o.Output)
	}
	// todo
	if cmd.Name() == "credential" {
		/*	cim := sp_api.NewCloudInfoManager()
			fmt.Println(cim)

			spidercli := app.Config.GetCurrentContext().SpiderCli

			err := cim.SetServerAddr(spidercli.ServerAddr)
			if err != nil {
				logger.Error("server_addr set failed", err)
			}

			timeout, _ := time.ParseDuration(spidercli.Timeout)
			err = cim.SetTimeout(timeout)
			if err != nil {
				logger.Error("timeout set failed", err)
			}
			err = cim.Open()
			if err != nil {
				logger.Error("spdier api open failed : ", err)
				return
			}
			defer cim.Close()
		*/
	}
	err = nil

	switch cmd.Parent().Name() {
	case "cbadm":
		switch cmd.Name() {
		case "healthy":
			result, err = mcar.Healthy()
		}
	case "get":
		switch cmd.Name() {
		case "cluster":
			if o.Name == "" {
				result, err = mcar.ListClusterByParam(o.Namespace)
			} else {
				result, err = mcar.GetClusterByParam(o.Namespace, o.Name)
			}
		case "node":
			if o.Name == "" {
				result, err = mcar.ListNodeByParam(o.Namespace, clusterName)
			} else {
				result, err = mcar.GetNodeByParam(o.Namespace, clusterName, o.Name)
			}
		case "credential":
			if o.Name == "" {
				//result, err = cim.ListCredential()
			} else {
				//result, err = cim.GetCredentialByParam(o.Name)
			}
		}
	case "create":
		switch cmd.Name() {
		case "cluster":
			result, err = mcar.CreateCluster(o.Data)
		case "node":
			result, err = mcar.AddNode(o.Data)
		case "credential":
			// result, err = cim.CreateCredential(o.Data)
		}
	case "delete":
		switch cmd.Name() {
		case "cluster":
			result, err = mcar.DeleteClusterByParam(o.Namespace, o.Name)
		case "node":
			result, err = mcar.RemoveNodeByParam(o.Namespace, clusterName, o.Name)
		case "credential":
			// result, err = cim.DeleteCredentialByParam(o.Name)
		}
	}

	if err != nil {
		if o.Output == "yaml" {
			fmt.Fprintf(cmd.OutOrStdout(), "message: %v\n", err)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "{\"message\": \"%v\"}\n", err)
		}
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", result)
	}
}
