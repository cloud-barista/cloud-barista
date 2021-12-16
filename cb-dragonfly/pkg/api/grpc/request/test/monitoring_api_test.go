package test

import (
	"fmt"
	"testing"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/protobuf/cbdragonfly"
	"github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/request"
)

func TestGRPC(t *testing.T) {
	// Test Connection
	monApi := request.InitMonitoringAPI()
	err := monApi.SetServerAddr("127.0.0.1:9999")
	if err != nil {
		t.Error(fmt.Sprintf("failed to set CB-Draognfly gRPC server config, error=%s", err))
	}
	err = monApi.Open()
	if err != nil {
		t.Error(fmt.Sprintf("failed to connect CB-Draognfly gRPC server, error=%s", err))
	}

	err = testMonConfig(monApi)
	if err != nil {
		t.Error(fmt.Sprintf("failed to test GET config API, error=%s", err))
	}
	err = testMonAgentInstall(monApi)
	if err != nil {
		t.Error(fmt.Sprintf("failed to test INSTALL agent API, error=%s", err))
	}
}

func testMonConfig(monApi *request.MonitoringAPI) error {
	result, err := monApi.GetMonitoringConfig()
	if err != nil {
		return err
	}
	fmt.Println("success to get mon config")
	fmt.Println(result)
	return nil
}

func testMonAgentInstall(monApi *request.MonitoringAPI) error {
	installAgentReq := cbdragonfly.InstallAgentRequest{
		NsId:     "test-ns",
		McisId:   "test-mcis",
		VmId:     "test-vm",
		PublicIp: "127.0.0.1",
		UserName: "cb-user",
		SshKey:   "-----BEGIN RSA PRIVATE KEY-----\n...",
		CspType:  "azure",
		Port:     "22",
	}
	result, err := monApi.InstallAgent(installAgentReq)
	if err != nil {
		return err
	}
	fmt.Println("success to install agent")
	fmt.Println(result)
	return nil
}
