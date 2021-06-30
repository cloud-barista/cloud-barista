// Proof of Concepts of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is a Cloud Driver Example for PoC Test.
//
// by hyokyung.kim@innogrid.co.kr, 2019.08.

package resources

import (
	"errors"
	"fmt"
	"strings"
	"time"

	call "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/call-log"
	"github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/drivers/cloudit/client"
	"github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/drivers/cloudit/client/ace/nic"
	"github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/drivers/cloudit/client/ace/server"
	"github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/drivers/cloudit/client/dna/adaptiveip"
	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
)

const (
	VMDefaultUser         = "root"
	VMDefaultPassword     = "qwe1212!Q"
	SSHDefaultUser        = "cb-user"
	SSHDefaultPort        = 22
	VM                    = "VM"
	DefaultSGName         = "ALL"
	ExtraInboundRuleName  = "extra-inbound"
	ExtraOutboundRuleName = "extra-outbound"
	InboundRule           = "inbound"
	OutboundRule          = "outbound"
)

type ClouditVMHandler struct {
	CredentialInfo idrv.CredentialInfo
	Client         *client.RestClient
}

func (vmHandler *ClouditVMHandler) StartVM(vmReqInfo irs.VMReqInfo) (irs.VMInfo, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vmHandler.Client.IdentityEndpoint, call.VM, VM, "StartVM()")

	// 가상서버 이름 중복 체크
	vmId, _ := vmHandler.getVmIdByName(vmReqInfo.IId.NameId)
	if vmId != "" {
		createErr := errors.New(fmt.Sprintf("VirtualMachine with name %s already exist", vmReqInfo.IId.NameId))
		LoggingError(hiscallInfo, createErr)
		return irs.VMInfo{}, createErr
	}

	// 이미지 정보 조회 (Name)
	imageHandler := ClouditImageHandler{
		Client:         vmHandler.Client,
		CredentialInfo: vmHandler.CredentialInfo,
	}
	image, err := imageHandler.GetImage(vmReqInfo.ImageIID)
	if err != nil {
		createErr := errors.New(fmt.Sprintf("failed to get image, err : %s", err.Error()))
		LoggingError(hiscallInfo, createErr)
		return irs.VMInfo{}, createErr
	}

	//  네트워크 정보 조회 (Name)
	vpcHandler := ClouditVPCHandler{
		Client:         vmHandler.Client,
		CredentialInfo: vmHandler.CredentialInfo,
	}
	vpc, err := vpcHandler.GetSubnet(vmReqInfo.SubnetIID)
	if err != nil {
		createErr := errors.New(fmt.Sprintf("failed to get virtual network, err : %s", err.Error()))
		LoggingError(hiscallInfo, createErr)
		return irs.VMInfo{}, createErr
	}

	// 보안그룹 정보 조회 (Name)
	sgHandler := ClouditSecurityHandler{
		Client:         vmHandler.Client,
		CredentialInfo: vmHandler.CredentialInfo,
	}

	// TODO: VNIC Security Group 수정 API 작업
	//// Default SG 조회
	//defaultSG, err := sgHandler.getSecurityByName(DefaultSGName)
	//if err != nil {
	//	return irs.VMInfo{}, nil
	//}
	//defaultSecGroups := make([]server.SecGroupInfo, len(vmReqInfo.SecurityGroupIIDs))
	//for i, _ := range vmReqInfo.SecurityGroupIIDs {
	//	defaultSecGroups[i] = server.SecGroupInfo{
	//		Id: defaultSG.ID,
	//	}
	//}

	// VM VNIC에 User Security Group 설정
	secGroups := make([]server.SecGroupInfo, len(vmReqInfo.SecurityGroupIIDs))
	for i, s := range vmReqInfo.SecurityGroupIIDs {
		secGroups[i] = server.SecGroupInfo{
			Id: s.SystemId,
		}
	}

	// 임시 Inbound 규칙 생성
	_, err = sgHandler.addRuleToSG(ExtraInboundRuleName, secGroups[0].Id, InboundRule)
	if err != nil {
		createErr := errors.New(fmt.Sprintf("failed to add extra inbound rule to SG, err : %s", err.Error()))
		LoggingError(hiscallInfo, createErr)
		return irs.VMInfo{}, createErr
	}
	// 임시 Outbound 규칙 생성
	_, err = sgHandler.addRuleToSG(ExtraOutboundRuleName, secGroups[0].Id, OutboundRule)
	if err != nil {
		createErr := errors.New(fmt.Sprintf("failed to add extra outbound rule to SG, err : %s", err.Error()))
		LoggingError(hiscallInfo, createErr)
		return irs.VMInfo{}, createErr
	}

	// Spec 정보 조회 (Name)
	vmSpecId, err := GetVMSpecByName(vmHandler.Client.AuthenticatedHeaders(), vmHandler.Client, vmReqInfo.VMSpecName)
	if err != nil {
		createErr := errors.New(fmt.Sprintf("failed to get vm spec, err : %s", err.Error()))
		LoggingError(hiscallInfo, createErr)
		return irs.VMInfo{}, createErr
	}

	vmHandler.Client.TokenID = vmHandler.CredentialInfo.AuthToken
	authHeader := vmHandler.Client.AuthenticatedHeaders()

	reqInfo := server.VMReqInfo{
		TemplateId:   image.IId.SystemId,
		SpecId:       *vmSpecId,
		Name:         vmReqInfo.IId.NameId,
		HostName:     vmReqInfo.IId.NameId,
		RootPassword: VMDefaultPassword,
		SubnetAddr:   vpc.Addr,
		Secgroups:    secGroups,
	}

	requestOpts := client.RequestOpts{
		MoreHeaders: authHeader,
		JSONBody:    reqInfo,
	}

	// VM 생성
	start := call.Start()
	creatingVm, err := server.Start(vmHandler.Client, &requestOpts)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.VMInfo{}, err
	}
	LoggingInfo(hiscallInfo, start)

	// VM 생성 완료까지 wait
	curRetryCnt := 0
	maxRetryCnt := 120
	for {
		// Check VM Deploy Status
		vmInfo, err := server.Get(vmHandler.Client, creatingVm.ID, &requestOpts)
		if err != nil {
			LoggingError(hiscallInfo, err)
			return irs.VMInfo{}, err
		}

		if vmInfo.PrivateIp != "" && getVmStatus(vmInfo.State) == irs.Running {
			ok, err := vmHandler.AssociatePublicIP(creatingVm.Name, vmInfo.PrivateIp)
			if !ok {
				LoggingError(hiscallInfo, err)
				return irs.VMInfo{}, err
			}
			break
		}
		time.Sleep(1 * time.Second)
		curRetryCnt++
		if curRetryCnt > maxRetryCnt {
			cblogger.Errorf(fmt.Sprintf("failed to start vm, exceeded maximum retry count %d", maxRetryCnt))
			return irs.VMInfo{}, errors.New(fmt.Sprintf("failed to start vm, exceeded maximum retry count %d", maxRetryCnt))
		}
	}

	vm, err := server.Get(vmHandler.Client, creatingVm.ID, &requestOpts)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.VMInfo{}, err
	}
	vmInfo := vmHandler.mappingServerInfo(*vm)

	// SSH 접속 사용자 및 공개키 등록
	loginUserId := SSHDefaultUser
	createUserErr := errors.New("Error adding cb-User to new VM")

	// SSH 접속까지 시도
	curConnectionCnt := 0
	maxConnectionRetryCnt := 30
	for {
		cblogger.Info("Trying to connect via root user ...")
		_, err := RunCommand(vmInfo.PublicIP, SSHDefaultPort, VMDefaultUser, VMDefaultPassword, "echo test")
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
		curConnectionCnt++
		if curConnectionCnt > maxConnectionRetryCnt {
			return irs.VMInfo{}, createUserErr
		}
	}

	// 사용자 등록 및 sudoer 권한 추가
	_, err = RunCommand(vmInfo.PublicIP, SSHDefaultPort, VMDefaultUser, VMDefaultPassword, fmt.Sprintf("useradd -s /bin/bash %s -rm", loginUserId))
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.VMInfo{}, createUserErr
	}
	_, err = RunCommand(vmInfo.PublicIP, SSHDefaultPort, VMDefaultUser, VMDefaultPassword, fmt.Sprintf("echo \"%s ALL=(root) NOPASSWD:ALL\" >> /etc/sudoers", loginUserId))
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.VMInfo{}, createUserErr
	}

	// 공개키 등록
	_, err = RunCommand(vmInfo.PublicIP, SSHDefaultPort, VMDefaultUser, VMDefaultPassword, fmt.Sprintf("mkdir -p /home/%s/.ssh", loginUserId))
	publicKey, err := GetPublicKey(vmHandler.CredentialInfo, vmReqInfo.KeyPairIID.NameId)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.VMInfo{}, createUserErr
	}
	_, err = RunCommand(vmInfo.PublicIP, SSHDefaultPort, VMDefaultUser, VMDefaultPassword, fmt.Sprintf("echo \"%s\" > /home/%s/.ssh/authorized_keys", publicKey, loginUserId))

	// ssh 접속 방법 변경 (sshd_config 파일 변경)
	_, err = RunCommand(vmInfo.PublicIP, SSHDefaultPort, VMDefaultUser, VMDefaultPassword, "sed -i 's/PasswordAuthentication yes/PasswordAuthentication no/g' /etc/ssh/sshd_config")
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.VMInfo{}, createUserErr
	}
	_, err = RunCommand(vmInfo.PublicIP, SSHDefaultPort, VMDefaultUser, VMDefaultPassword, "sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/g' /etc/ssh/sshd_config")
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.VMInfo{}, createUserErr
	}
	_, err = RunCommand(vmInfo.PublicIP, SSHDefaultPort, VMDefaultUser, VMDefaultPassword, "systemctl restart sshd")
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.VMInfo{}, createUserErr
	}

	// VM VNIC에 User Security Group Attach
	userSecGroups := make([]server.SecGroupInfo, len(vmReqInfo.SecurityGroupIIDs))
	for i, s := range vmReqInfo.SecurityGroupIIDs {
		userSecGroups[i] = server.SecGroupInfo{
			Id: s.SystemId,
		}
	}

	// SG 임시 규칙 삭제를 위한 Rule ID 조회
	extraSG, _ := sgHandler.listRulesInSG(secGroups[0].Id)
	//deleteTargetRuleID := make([]string, 2)
	var deleteTargetRuleID []string

	for _, v := range *extraSG {
		if v.Name == ExtraInboundRuleName || v.Name == ExtraOutboundRuleName {
			deleteTargetRuleID = append(deleteTargetRuleID, v.ID)
		}
	}

	// SG 임시 규칙 삭제
	for _, v := range deleteTargetRuleID {
		err = sgHandler.deleteRuleInSG(secGroups[0].Id, v)
		if err != nil {
			deleteErr := errors.New(fmt.Sprintf("failed to delete extra rules, err : %s", err.Error()))
			LoggingError(hiscallInfo, deleteErr)
			return irs.VMInfo{}, deleteErr
		}
	}
	// TODO: VNIC의 SG 설정 API 수정
	//err = vmHandler.attachSgToVnic(authHeader, vm.ID, vmHandler.Client, vnicMac, defaultSecGroups, userSecGroups)
	//if err != nil {
	//	return irs.VMInfo{}, err
	//}

	return vmInfo, nil
}

func (vmHandler *ClouditVMHandler) SuspendVM(vmIID irs.IID) (irs.VMStatus, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vmHandler.Client.IdentityEndpoint, call.VM, vmIID.NameId, "SuspendVM()")

	vmHandler.Client.TokenID = vmHandler.CredentialInfo.AuthToken
	authHeader := vmHandler.Client.AuthenticatedHeaders()

	requestOpts := client.RequestOpts{
		MoreHeaders: authHeader,
	}

	start := call.Start()
	err := server.Suspend(vmHandler.Client, vmIID.SystemId, &requestOpts)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.Failed, err
	}
	LoggingInfo(hiscallInfo, start)

	// VM 상태 정보 반환
	vmStatus, err := vmHandler.GetVMStatus(vmIID)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.Failed, err
	}
	return vmStatus, nil
}

func (vmHandler *ClouditVMHandler) ResumeVM(vmIID irs.IID) (irs.VMStatus, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vmHandler.Client.IdentityEndpoint, call.VM, vmIID.NameId, "ResumeVM()")

	vmHandler.Client.TokenID = vmHandler.CredentialInfo.AuthToken
	authHeader := vmHandler.Client.AuthenticatedHeaders()

	requestOpts := client.RequestOpts{
		MoreHeaders: authHeader,
	}

	start := call.Start()
	err := server.Resume(vmHandler.Client, vmIID.SystemId, &requestOpts)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.Failed, err
	}
	LoggingInfo(hiscallInfo, start)

	// VM 상태 정보 반환
	vmStatus, err := vmHandler.GetVMStatus(vmIID)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.Failed, err
	}
	return vmStatus, nil
}

func (vmHandler *ClouditVMHandler) RebootVM(vmIID irs.IID) (irs.VMStatus, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vmHandler.Client.IdentityEndpoint, call.VM, vmIID.NameId, "RebootVM()")

	vmHandler.Client.TokenID = vmHandler.CredentialInfo.AuthToken
	authHeader := vmHandler.Client.AuthenticatedHeaders()

	requestOpts := client.RequestOpts{
		MoreHeaders: authHeader,
	}

	start := call.Start()
	err := server.Reboot(vmHandler.Client, vmIID.SystemId, &requestOpts)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.Failed, err
	}
	LoggingInfo(hiscallInfo, start)

	// VM 상태 정보 반환
	vmStatus, err := vmHandler.GetVMStatus(vmIID)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.Failed, err
	}
	return vmStatus, nil
}

func (vmHandler *ClouditVMHandler) TerminateVM(vmIID irs.IID) (irs.VMStatus, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vmHandler.Client.IdentityEndpoint, call.VM, vmIID.NameId, "TerminateVM()")

	vmHandler.Client.TokenID = vmHandler.CredentialInfo.AuthToken
	authHeader := vmHandler.Client.AuthenticatedHeaders()

	requestOpts := client.RequestOpts{
		MoreHeaders: authHeader,
	}

	// VM 정보 조회
	vmInfo, err := vmHandler.GetVM(vmIID)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.Failed, err
	}

	// 연결된 PublicIP 반환
	if vmInfo.PublicIP != "" {
		if ok, err := vmHandler.DisassociatePublicIP(vmInfo.PublicIP); !ok {
			LoggingError(hiscallInfo, err)
			return irs.Failed, err
		}
		time.Sleep(5 * time.Second)
	}

	start := call.Start()
	if err := server.Terminate(vmHandler.Client, vmInfo.IId.SystemId, &requestOpts); err != nil {
		LoggingError(hiscallInfo, err)
		return irs.Failed, err
	}
	LoggingInfo(hiscallInfo, start)

	// VM 상태 정보 반환
	return irs.Terminating, nil
}

func (vmHandler *ClouditVMHandler) ListVMStatus() ([]*irs.VMStatusInfo, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vmHandler.Client.IdentityEndpoint, call.VM, VM, "ListVMStatus()")

	vmHandler.Client.TokenID = vmHandler.CredentialInfo.AuthToken
	authHeader := vmHandler.Client.AuthenticatedHeaders()

	requestOpts := client.RequestOpts{
		MoreHeaders: authHeader,
	}

	start := call.Start()
	vmList, err := server.List(vmHandler.Client, &requestOpts)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return []*irs.VMStatusInfo{}, err
	}
	LoggingInfo(hiscallInfo, start)

	vmStatusList := make([]*irs.VMStatusInfo, len(*vmList))
	for i, vm := range *vmList {
		vmStatusInfo := irs.VMStatusInfo{
			IId: irs.IID{
				NameId:   vm.Name,
				SystemId: vm.ID,
			},
			VmStatus: irs.VMStatus(vm.State),
		}
		vmStatusList[i] = &vmStatusInfo
	}
	return vmStatusList, nil
}

func (vmHandler *ClouditVMHandler) GetVMStatus(vmIID irs.IID) (irs.VMStatus, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vmHandler.Client.IdentityEndpoint, call.VM, vmIID.NameId, "GetVMStatus()")

	vmHandler.Client.TokenID = vmHandler.CredentialInfo.AuthToken
	authHeader := vmHandler.Client.AuthenticatedHeaders()

	requestOpts := client.RequestOpts{
		MoreHeaders: authHeader,
	}

	vmSystemID, err := vmHandler.getVmIdByName(vmIID.NameId)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.Failed, err
	}

	start := call.Start()
	vm, err := server.Get(vmHandler.Client, vmSystemID, &requestOpts)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.Failed, err
	}
	LoggingInfo(hiscallInfo, start)

	// Set VM Status Info
	status := getVmStatus(vm.State)
	return status, nil
}

func (vmHandler *ClouditVMHandler) ListVM() ([]*irs.VMInfo, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vmHandler.Client.IdentityEndpoint, call.VM, VM, "ListVM()")

	vmHandler.Client.TokenID = vmHandler.CredentialInfo.AuthToken
	authHeader := vmHandler.Client.AuthenticatedHeaders()

	requestOpts := client.RequestOpts{
		MoreHeaders: authHeader,
	}

	start := call.Start()
	vmList, err := server.List(vmHandler.Client, &requestOpts)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return []*irs.VMInfo{}, err
	}
	LoggingInfo(hiscallInfo, start)

	vmInfoList := make([]*irs.VMInfo, len(*vmList))
	for i, vm := range *vmList {
		vmInfo := vmHandler.mappingServerInfo(vm)
		vmInfoList[i] = &vmInfo
	}
	return vmInfoList, nil
}

func (vmHandler *ClouditVMHandler) GetVM(vmIID irs.IID) (irs.VMInfo, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vmHandler.Client.IdentityEndpoint, call.VM, vmIID.NameId, "GetVM()")

	vmHandler.Client.TokenID = vmHandler.CredentialInfo.AuthToken
	authHeader := vmHandler.Client.AuthenticatedHeaders()

	requestOpts := client.RequestOpts{
		MoreHeaders: authHeader,
	}

	start := call.Start()
	vm, err := server.Get(vmHandler.Client, vmIID.SystemId, &requestOpts)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.VMInfo{}, err
	}
	LoggingInfo(hiscallInfo, start)

	vmInfo := vmHandler.mappingServerInfo(*vm)
	return vmInfo, nil
}

// VM에 PublicIP 연결
func (vmHandler *ClouditVMHandler) AssociatePublicIP(vmName string, vmIp string) (bool, error) {
	vmHandler.Client.TokenID = vmHandler.CredentialInfo.AuthToken
	authHeader := vmHandler.Client.AuthenticatedHeaders()

	var availableIP adaptiveip.IPInfo

	// 1. 사용 가능한 PublicIP 목록 가져오기
	requestOpts := client.RequestOpts{
		MoreHeaders: authHeader,
	}
	if availableIPList, err := adaptiveip.ListAvailableIP(vmHandler.Client, &requestOpts); err != nil {
		return false, err
	} else {
		if len(*availableIPList) == 0 {
			allocateErr := errors.New(fmt.Sprintf("There is no PublicIPs to allocate"))
			return false, allocateErr
		} else {
			availableIP = (*availableIPList)[0]
		}
	}

	// 2. PublicIP 생성 및 할당
	reqInfo := adaptiveip.PublicIPReqInfo{
		IP:        availableIP.IP,
		Name:      vmName + "-PublicIP",
		PrivateIP: vmIp,
	}

	createOpts := client.RequestOpts{
		JSONBody:    reqInfo,
		MoreHeaders: authHeader,
	}
	_, err := adaptiveip.Create(vmHandler.Client, &createOpts)
	if err != nil {
		cblogger.Error(err)
		return false, err
	}
	return true, nil
}

// VM에 PublicIP 해제
func (vmHandler *ClouditVMHandler) DisassociatePublicIP(publicIP string) (bool, error) {
	vmHandler.Client.TokenID = vmHandler.CredentialInfo.AuthToken
	authHeader := vmHandler.Client.AuthenticatedHeaders()

	requestOpts := client.RequestOpts{
		MoreHeaders: authHeader,
	}

	if err := adaptiveip.Delete(vmHandler.Client, publicIP, &requestOpts); err != nil {
		cblogger.Error(err)
		return false, err
	} else {
		return true, nil
	}
}

func (vmHandler *ClouditVMHandler) mappingServerInfo(server server.ServerInfo) irs.VMInfo {

	// Get Default VM Info
	vmInfo := irs.VMInfo{
		IId: irs.IID{
			NameId:   server.Name,
			SystemId: server.ID,
		},
		Region: irs.RegionInfo{
			Region: server.TenantID,
			Zone:   server.TenantID,
		},
		ImageIId: irs.IID{
			NameId:   server.Template,
			SystemId: server.TemplateID,
		},
		VMSpecName: server.Spec,
		VpcIID: irs.IID{
			NameId:   defaultVPCName,
			SystemId: defaultVPCName,
		},
		VMUserId:  VMDefaultUser,
		PublicIP:  server.AdaptiveIp,
		PrivateIP: server.PrivateIp,
	}

	if server.CreatedAt != "" {
		timeArr := strings.Split(server.CreatedAt, " ")
		timeFormatStr := fmt.Sprintf("%sT%sZ", timeArr[0], timeArr[1])
		if createTime, err := time.Parse(time.RFC3339, timeFormatStr); err == nil {
			vmInfo.StartTime = createTime
		}
	}

	// Get Subnet Info
	VPCHandler := ClouditVPCHandler{
		Client:         vmHandler.Client,
		CredentialInfo: vmHandler.CredentialInfo,
	}
	subnet, err := VPCHandler.GetSubnet(irs.IID{NameId: server.SubnetAddr})
	if err == nil {
		vmInfo.SubnetIID = irs.IID{
			NameId:   subnet.Name,
			SystemId: subnet.ID,
		}
	}

	// Get SecurityGroup Info
	vmHandler.Client.TokenID = vmHandler.CredentialInfo.AuthToken
	authHeader := vmHandler.Client.AuthenticatedHeaders()
	vnicList, _ := ListVNic(authHeader, vmHandler.Client, server.ID)
	if vnicList != nil {
		defaultVnic := (*vnicList)[0]
		segGroupList := make([]irs.IID, len(defaultVnic.SecGroups))
		for i, s := range defaultVnic.SecGroups {
			segGroupList[i] = irs.IID{
				NameId:   s.Name,
				SystemId: s.Id,
			}
		}
		vmInfo.SecurityGroupIIds = segGroupList
	}
	return vmInfo
}

func (vmHandler *ClouditVMHandler) getVmIdByName(vmNameID string) (string, error) {
	var vmId string

	// VM 목록 검색
	vmList, err := vmHandler.ListVM()
	if err != nil {
		return "", err
	}

	// VM 목록에서 Name 기준 검색
	for _, v := range vmList {
		if strings.EqualFold(v.IId.NameId, vmNameID) {
			vmId = v.IId.SystemId
			break
		}
	}

	// 만약 VM이 검색되지 않을 경우 에러 처리
	if vmId == "" {
		err := errors.New(fmt.Sprintf("failed to find vm with name %s", vmNameID))
		return "", err
	}
	return vmId, nil
}

func getVmStatus(vmStatus string) irs.VMStatus {
	var resultStatus string
	switch strings.ToLower(vmStatus) {
	case "creating":
		resultStatus = "Creating"
	case "running":
		resultStatus = "Running"
	case "stopping":
		resultStatus = "Suspending"
	case "stopped":
		resultStatus = "Suspended"
	case "starting":
		resultStatus = "Resuming"
	case "rebooting":
		resultStatus = "Rebooting"
	case "terminating":
		resultStatus = "Terminating"
	case "terminated":
		resultStatus = "Terminated"
	case "failed":
	default:
		resultStatus = "Failed"
	}
	return irs.VMStatus(resultStatus)
}

func (vmHandler *ClouditVMHandler) attachSgToVnic(authHeader map[string]string, vmID string, reqClient *client.RestClient, vnicMac string, sgGroup []server.SecGroupInfo) {

	reqInfo := server.VMReqInfo{
		Secgroups: sgGroup,
	}
	requestOpts := client.RequestOpts{
		MoreHeaders: authHeader,
		JSONBody:    reqInfo,
	}
	nic.Put(reqClient, vmID, &requestOpts, vnicMac)
}
