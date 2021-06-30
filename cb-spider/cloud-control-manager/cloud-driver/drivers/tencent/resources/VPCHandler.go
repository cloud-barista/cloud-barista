// Cloud Driver Interface of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is Resouces interfaces of Cloud Driver.
//
// by devunet@mz.co.kr

package resources

import (
	"errors"
	"strconv"

	call "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/call-log"
	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
	"github.com/davecgh/go-spew/spew"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

type TencentVPCHandler struct {
	Region idrv.RegionInfo
	Client *vpc.Client
}

func (VPCHandler *TencentVPCHandler) CreateVPC(vpcReqInfo irs.VPCReqInfo) (irs.VPCInfo, error) {
	cblogger.Info(vpcReqInfo)

	zoneId := VPCHandler.Region.Zone
	cblogger.Infof("Zone : %s", zoneId)
	if zoneId == "" {
		cblogger.Error("Connection 정보에 Zone 정보가 없습니다.")
		return irs.VPCInfo{}, errors.New("Connection 정보에 Zone 정보가 없습니다.")
	}

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.TENCENT,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: vpcReqInfo.IId.NameId,
		CloudOSAPI:   "CreateVpc()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}

	//=========================
	// VPC 생성
	//=========================
	request := vpc.NewCreateVpcRequest()
	request.VpcName = common.StringPtr(vpcReqInfo.IId.NameId)
	request.CidrBlock = common.StringPtr(vpcReqInfo.IPv4_CIDR)

	callLogStart := call.Start()
	response, err := VPCHandler.Client.CreateVpc(request)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	cblogger.Debug(response.ToJsonString())
	//spew.Dump(result)
	if err != nil {
		callLogInfo.ErrorMSG = err.Error()
		callogger.Error(call.String(callLogInfo))
		cblogger.Error(err)
		return irs.VPCInfo{}, err
	}
	callogger.Info(call.String(callLogInfo))

	newVpcId := *response.Response.Vpc.VpcId // Subnet이 포함된 정보를 전달해야 하기 때문에 생성된 VPC Id를 보관함.

	//=========================
	// Subnet 생성
	//========================
	requestSubnet := vpc.NewCreateSubnetsRequest()

	requestSubnet.VpcId = common.StringPtr(newVpcId)
	requestSubnet.Subnets = []*vpc.SubnetInput{}

	for _, curSubnet := range vpcReqInfo.SubnetInfoList {
		cblogger.Infof("[%s] Subnet 처리", curSubnet.IId.NameId)
		reqSubnet := &vpc.SubnetInput{
			CidrBlock:  common.StringPtr(curSubnet.IPv4_CIDR),
			SubnetName: common.StringPtr(curSubnet.IId.NameId),
			Zone:       common.StringPtr(zoneId),
			//RouteTableId: common.StringPtr("route"),
		}
		requestSubnet.Subnets = append(requestSubnet.Subnets, reqSubnet)
	}

	responseSubnet, errSubnet := VPCHandler.Client.CreateSubnets(requestSubnet)
	cblogger.Debug(responseSubnet.ToJsonString())
	spew.Dump(responseSubnet)
	if errSubnet != nil {
		cblogger.Error(errSubnet)
		return irs.VPCInfo{}, errSubnet
	}

	//신규로 생성된 VPC와 Subnet 정보를 irs.VPCInfo{}로 치환해도 되지만 수정의 편의및 최신 정보 통일을 위해 GetVPC롤 호출함.
	//생성된 Subnet을 포함한 VPC의 최신 정보를 조회함.
	retVpcInfo, errVpc := VPCHandler.GetVPC(irs.IID{SystemId: newVpcId})
	if errVpc != nil {
		cblogger.Error(errVpc)
		return irs.VPCInfo{}, errVpc
	}
	retVpcInfo.IId.NameId = vpcReqInfo.IId.NameId // 생성 시에는 NameId는 cb-spider를 위해 요청 받은 값을 그대로 리턴해야 함.

	return retVpcInfo, nil
}

//VPC 정보를 추출함
func ExtractVpcDescribeInfo(vpcInfo *vpc.Vpc) irs.VPCInfo {
	// cblogger.Debug("전달 받은 내용")
	// spew.Dump(vpcInfo)
	resVpcInfo := irs.VPCInfo{
		//NameId는 사용되지 않기 때문에 전달할 필요가 없지만 Tencent는 Name도 필수로 들어가니 전달함.
		IId:       irs.IID{SystemId: *vpcInfo.VpcId, NameId: *vpcInfo.VpcName},
		IPv4_CIDR: *vpcInfo.CidrBlock,
	}

	return resVpcInfo
}

func (VPCHandler *TencentVPCHandler) ListVPC() ([]*irs.VPCInfo, error) {
	cblogger.Info("Start")

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.TENCENT,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: "ListVPC",
		CloudOSAPI:   "DescribeVpcs()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}

	request := vpc.NewDescribeVpcsRequest()
	callLogStart := call.Start()
	response, err := VPCHandler.Client.DescribeVpcs(request)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	cblogger.Debug(response.ToJsonString())
	//spew.Dump(result)
	if err != nil {
		callLogInfo.ErrorMSG = err.Error()
		callogger.Error(call.String(callLogInfo))
		cblogger.Error(err)
		return nil, err
	}
	callogger.Info(call.String(callLogInfo))

	cblogger.Info("VPC 개수 : ", *response.Response.TotalCount)

	var vpcInfoList []*irs.VPCInfo
	if *response.Response.TotalCount > 0 {
		for _, curVpc := range response.Response.VpcSet {
			cblogger.Debugf("[%s] VPC 정보 조회 - [%s]", *curVpc.VpcId, *curVpc.VpcName)
			vpcInfo, vpcErr := VPCHandler.GetVPC(irs.IID{SystemId: *curVpc.VpcId})
			// cblogger.Info("==>조회 결과")
			// spew.Dump(vpcInfo)
			if vpcErr != nil {
				cblogger.Error(vpcErr)
				return nil, vpcErr
			}
			vpcInfoList = append(vpcInfoList, &vpcInfo)
		}
	}

	cblogger.Debugf("리턴 결과 목록 수 : [%d]", len(vpcInfoList))
	// spew.Dump(vpcInfoList)
	return vpcInfoList, nil
}

func (VPCHandler *TencentVPCHandler) GetVPC(vpcIID irs.IID) (irs.VPCInfo, error) {
	cblogger.Info("VPC IID : ", vpcIID.SystemId)

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.TENCENT,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: "GetVPC",
		CloudOSAPI:   "DescribeVpcs()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}

	request := vpc.NewDescribeVpcsRequest()
	request.VpcIds = common.StringPtrs([]string{vpcIID.SystemId})

	callLogStart := call.Start()
	response, err := VPCHandler.Client.DescribeVpcs(request)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	if err != nil {
		cblogger.Errorf("An API error has returned: %s", err.Error())
		callLogInfo.ErrorMSG = err.Error()
		callogger.Error(call.String(callLogInfo))
		return irs.VPCInfo{}, err
	}
	callogger.Info(call.String(callLogInfo))

	cblogger.Debug("VPC 개수 : ", *response.Response.TotalCount)
	if *response.Response.TotalCount < 1 {
		return irs.VPCInfo{}, errors.New("Notfound: '" + vpcIID.SystemId + "' VPC Not found")
	}

	vpcInfo := ExtractVpcDescribeInfo(response.Response.VpcSet[0])
	cblogger.Debug(vpcInfo)

	//=======================
	// Subnet 처리
	//=======================
	var errSubnet error
	vpcInfo.SubnetInfoList, errSubnet = VPCHandler.ListSubnet(vpcIID.SystemId)
	if errSubnet != nil {
		callogger.Error(errSubnet)
		return vpcInfo, errSubnet
	}

	return vpcInfo, nil
}

func (VPCHandler *TencentVPCHandler) DeleteVPC(vpcIID irs.IID) (bool, error) {
	cblogger.Infof("Delete VPC : [%s]", vpcIID.SystemId)

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.TENCENT,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: vpcIID.SystemId,
		CloudOSAPI:   "DeleteVpc()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}

	request := vpc.NewDeleteVpcRequest()
	request.VpcId = common.StringPtr(vpcIID.SystemId)

	callLogStart := call.Start()
	_, err := VPCHandler.Client.DeleteVpc(request)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	if err != nil {
		cblogger.Errorf("An API error has returned: %s", err.Error())
		callLogInfo.ErrorMSG = err.Error()
		callogger.Error(call.String(callLogInfo))
		return false, err
	}
	callogger.Info(call.String(callLogInfo))

	return true, nil
}

func (VPCHandler *TencentVPCHandler) ListSubnet(reqVpcId string) ([]irs.SubnetInfo, error) {
	cblogger.Infof("reqVpcId : [%s]", reqVpcId)
	var arrSubnetInfoList []irs.SubnetInfo

	/*
		// logger for HisCall
		callogger := call.GetLogger("HISCALL")
		callLogInfo := call.CLOUDLOGSCHEMA{
			CloudOS:      call.TENCENT,
			RegionZone:   VPCHandler.Region.Zone,
			ResourceType: call.VPCSUBNET,
			ResourceName: "ListSubnet - VpcId:" + reqVpcId,
			CloudOSAPI:   "DescribeSubnets()",
			ElapsedTime:  "",
			ErrorMSG:     "",
		}
	*/

	request := vpc.NewDescribeSubnetsRequest()
	request.Filters = []*vpc.Filter{
		&vpc.Filter{
			Name:   common.StringPtr("vpc-id"),
			Values: common.StringPtrs([]string{reqVpcId}),
		},
	}

	// callLogStart := call.Start()
	response, err := VPCHandler.Client.DescribeSubnets(request)
	// callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	//cblogger.Debug(response.ToJsonString())
	spew.Dump(response)
	if err != nil {
		// callLogInfo.ErrorMSG = err.Error()
		// callogger.Error(call.String(callLogInfo))
		cblogger.Error(err)
		return nil, err
	}
	// callogger.Info(call.String(callLogInfo))

	for _, curSubnet := range response.Response.SubnetSet {
		cblogger.Infof("[%s] Subnet 정보 조회", *curSubnet.SubnetId)
		resSubnetInfo := irs.SubnetInfo{
			IId:       irs.IID{SystemId: *curSubnet.SubnetId, NameId: *curSubnet.SubnetName},
			IPv4_CIDR: *curSubnet.CidrBlock,
			//Status:    *subnetInfo.State,
		}

		keyValueList := []irs.KeyValue{
			{Key: "VpcId", Value: *curSubnet.VpcId},
			{Key: "IsDefault", Value: strconv.FormatBool(*curSubnet.IsDefault)},
			{Key: "AvailabilityZone", Value: *curSubnet.Zone},
		}
		resSubnetInfo.KeyValueList = keyValueList
		arrSubnetInfoList = append(arrSubnetInfoList, resSubnetInfo)
	}

	return arrSubnetInfoList, nil
}

func (VPCHandler *TencentVPCHandler) isExistSubnet(reqSubnetId string) (bool, error) {
	cblogger.Infof("reqSubnetId : [%s]", reqSubnetId)

	request := vpc.NewDescribeSubnetsRequest()
	request.Filters = []*vpc.Filter{
		&vpc.Filter{
			Name:   common.StringPtr("subnet-name"),
			Values: common.StringPtrs([]string{reqSubnetId}),
		},
	}

	//spew.Dump(request)
	response, err := VPCHandler.Client.DescribeSubnets(request)
	cblogger.Info("서브넷 실행 결과")
	//spew.Dump(response)
	if err != nil {
		cblogger.Error(err)
		return false, err
	}

	if *response.Response.TotalCount < 1 {
		return false, nil
	}

	return true, nil
}

func (VPCHandler *TencentVPCHandler) AddSubnet(vpcIID irs.IID, subnetInfo irs.SubnetInfo) (irs.VPCInfo, error) {
	cblogger.Infof("[%s] Subnet 추가 - CIDR : %s", subnetInfo.IId.NameId, subnetInfo.IPv4_CIDR)

	zoneId := VPCHandler.Region.Zone
	cblogger.Infof("Zone : %s", zoneId)
	if zoneId == "" {
		cblogger.Error("Connection 정보에 Zone 정보가 없습니다.")
		return irs.VPCInfo{}, errors.New("Connection 정보에 Zone 정보가 없습니다.")
	}

	if subnetInfo.IId.NameId == "" {
		return irs.VPCInfo{}, errors.New("생성할 SubnetId 정보가 없습니다.")
	}

	isExit, errSubnetInfo := VPCHandler.isExistSubnet(subnetInfo.IId.NameId)
	if errSubnetInfo != nil {
		cblogger.Error(errSubnetInfo)
		return irs.VPCInfo{}, errSubnetInfo
	}

	cblogger.Info("Subnet 존재여부 : ")
	cblogger.Info(isExit)

	if isExit {
		cblogger.Errorf("이미 [%S] Subnet이 존재하기 때문에 생성하지 않고 기존 정보와 함께 에러를 리턴함.", subnetInfo.IId.NameId)
		return irs.VPCInfo{}, errors.New("InvalidVNetwork.Duplicate: The Subnet '" + subnetInfo.IId.NameId + "' already exists.")
	}

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.TENCENT,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: vpcIID.SystemId,
		CloudOSAPI:   "CreateSubnet()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}

	request := vpc.NewCreateSubnetRequest()

	request.VpcId = common.StringPtr(vpcIID.SystemId)
	request.SubnetName = common.StringPtr(subnetInfo.IId.NameId)
	request.CidrBlock = common.StringPtr(subnetInfo.IPv4_CIDR)
	request.Zone = common.StringPtr(VPCHandler.Region.Zone)

	callLogStart := call.Start()
	response, err := VPCHandler.Client.CreateSubnet(request)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	cblogger.Debug(response.ToJsonString())
	spew.Dump(response)
	if err != nil {
		callLogInfo.ErrorMSG = err.Error()
		callogger.Error(call.String(callLogInfo))
		cblogger.Error(err)
		return irs.VPCInfo{}, err
	}
	callogger.Info(call.String(callLogInfo))

	retVpcInfo, errVpcInfo := VPCHandler.GetVPC(vpcIID)
	if errVpcInfo != nil {
		cblogger.Error(errVpcInfo)
		return irs.VPCInfo{}, err
	}

	//retVpcInfo.SubnetInfoList[0].IId.NameId = vpcReqInfo.IId.NameId // 생성 시에는 NameId는 요청 받은 값으로 리턴해야 함.

	return retVpcInfo, nil
}

func (VPCHandler *TencentVPCHandler) RemoveSubnet(vpcIID irs.IID, subnetIID irs.IID) (bool, error) {
	cblogger.Infof("[%s] VPC의 [%s] Subnet 삭제", vpcIID.SystemId, subnetIID.SystemId)

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.TENCENT,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: vpcIID.SystemId,
		CloudOSAPI:   "DeleteSubnet()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}

	request := vpc.NewDeleteSubnetRequest()
	request.SubnetId = common.StringPtr(subnetIID.SystemId)

	callLogStart := call.Start()
	response, err := VPCHandler.Client.DeleteSubnet(request)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	cblogger.Debug(response.ToJsonString())
	//spew.Dump(response)
	if err != nil {
		callLogInfo.ErrorMSG = err.Error()
		callogger.Error(call.String(callLogInfo))
		cblogger.Error(err)
		return false, err
	}
	callogger.Info(call.String(callLogInfo))

	return true, nil
}

/*
func (VPCHandler *TencentVPCHandler) CreateVPC(vpcReqInfo irs.VPCReqInfo) (irs.VPCInfo, error) {
	cblogger.Info(vpcReqInfo)

	zoneId := VPCHandler.Region.Zone
	cblogger.Infof("Zone : %s", zoneId)
	if zoneId == "" {
		cblogger.Error("Connection 정보에 Zone 정보가 없습니다.")
		return irs.VPCInfo{}, errors.New("Connection 정보에 Zone 정보가 없습니다.")
	}

	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(vpcReqInfo.IPv4_CIDR),
	}

	spew.Dump(input)
	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.AWS,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: vpcReqInfo.IId.NameId,
		CloudOSAPI:   "CreateVpc()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	callLogStart := call.Start()

	result, err := VPCHandler.Client.CreateVpc(input)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
				callLogInfo.ErrorMSG = aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
			callLogInfo.ErrorMSG = err.Error()
		}
		callogger.Info(call.String(callLogInfo))
		return irs.VPCInfo{}, err
	}
	callogger.Info(call.String(callLogInfo))

	cblogger.Info(result)
	spew.Dump(result)
	retVpcInfo := ExtractVpcDescribeInfo(result.Vpc)
	retVpcInfo.IId.NameId = vpcReqInfo.IId.NameId // NameId는 요청 받은 값으로 리턴해야 함.

	//IGW Name Tag 설정
	if SetNameTag(VPCHandler.Client, *result.Vpc.VpcId, vpcReqInfo.IId.NameId) {
		cblogger.Infof("VPC에 %s Name 설정 성공", vpcReqInfo.IId.NameId)
	} else {
		cblogger.Errorf("VPC에 %s Name 설정 실패", vpcReqInfo.IId.NameId)
	}

	//====================================
	// PublicIP 할당을 위해 IGW 생성및 연결
	//====================================
	//IGW 생성
	resultIGW, errIGW := VPCHandler.Client.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})
	if errIGW != nil {
		if aerr, ok := errIGW.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(errIGW.Error())
		}
		return retVpcInfo, errIGW
	}

	cblogger.Info(resultIGW)

	//IGW Name Tag 설정
	if SetNameTag(VPCHandler.Client, *resultIGW.InternetGateway.InternetGatewayId, vpcReqInfo.IId.NameId) {
		cblogger.Infof("IGW에 %s Name 설정 성공", vpcReqInfo.IId.NameId)
	} else {
		cblogger.Errorf("IGW에 %s Name 설정 실패", vpcReqInfo.IId.NameId)
	}

	// VPC에 IGW연결
	inputIGW := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(*resultIGW.InternetGateway.InternetGatewayId),
		VpcId:             aws.String(retVpcInfo.IId.SystemId),
	}

	resultIGWAttach, errIGWAttach := VPCHandler.Client.AttachInternetGateway(inputIGW)
	if err != nil {
		if aerr, ok := errIGWAttach.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(errIGWAttach.Error())
		}
		return retVpcInfo, errIGWAttach
	}

	cblogger.Info(resultIGWAttach)

	// 생성된 VPC의 기본 라우팅 테이블에 IGW 라우팅 정보 추가
	errRoute := VPCHandler.CreateRouteIGW(retVpcInfo.IId.SystemId, *resultIGW.InternetGateway.InternetGatewayId)
	if errRoute != nil {
		return retVpcInfo, errRoute
	}

	//==========================
	// Subnet 생성
	//==========================
	//VPCHandler.CreateSubnet(retVpcInfo.IId.SystemId, vpcReqInfo.SubnetInfoList[0])
	var resSubnetList []irs.SubnetInfo
	for _, curSubnet := range vpcReqInfo.SubnetInfoList {
		cblogger.Infof("[%s] Subnet 생성", curSubnet.IId.NameId)
		cblogger.Infof("Reqt Subnet Info [%v]", curSubnet)
		resSubnet, errSubnet := VPCHandler.CreateSubnet(retVpcInfo.IId.SystemId, curSubnet)

		if errSubnet != nil {
			return retVpcInfo, errSubnet
		}
		resSubnetList = append(resSubnetList, resSubnet)
	}
	retVpcInfo.SubnetInfoList = resSubnetList
	return retVpcInfo, nil
}

// 생성된 VPC의 라우팅 테이블에 IGW(Internet Gateway) 라우팅 정보를 생성함 (AWS 콘솔의 라우팅 테이블의 [라우팅] Tab 처리)
func (VPCHandler *TencentVPCHandler) CreateRouteIGW(vpcId string, igwId string) error {
	cblogger.Infof("VPC ID : [%s] / IGW ID : [%s]", vpcId, igwId)
	routeTableId, errRoute := VPCHandler.GetDefaultRouteTable(vpcId)
	if errRoute != nil {
		return errRoute
	}

	cblogger.Infof("RouteTable[%s]에 IGW[%s]에 대한 라우팅(0.0.0.0/0) 정보를 추가 합니다.", routeTableId, igwId)
	input := &ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            aws.String(igwId),
		RouteTableId:         aws.String(routeTableId),
	}

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")

	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.AWS,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: igwId,
		CloudOSAPI:   "CreateRoute()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	callLogStart := call.Start()

	result, err := VPCHandler.Client.CreateRoute(input)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	if err != nil {
		cblogger.Errorf("RouteTable[%s]에 IGW[%s]에 대한 라우팅(0.0.0.0/0) 정보 추가 실패", routeTableId, igwId)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
				callLogInfo.ErrorMSG = aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
			callLogInfo.ErrorMSG = err.Error()
		}
		callogger.Info(call.String(callLogInfo))
		return err
	}
	cblogger.Infof("RouteTable[%s]에 IGW[%s]에 대한 라우팅(0.0.0.0/0) 정보를 추가 완료", routeTableId, igwId)
	callogger.Info(call.String(callLogInfo))

	cblogger.Info(result)
	spew.Dump(result)
	return nil
}

//https://docs.aws.amazon.com/ko_kr/vpc/latest/userguide/VPC_Route_Tables.html
//https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeRouteTables.html
// 자동 생성된 VPC의 기본 라우팅 테이블 정보를 찾음
func (VPCHandler *TencentVPCHandler) GetDefaultRouteTable(vpcId string) (string, error) {
	input := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(vpcId),
				},
			},
		},
	}

	result, err := VPCHandler.Client.DescribeRouteTables(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
		}
		return "", err
	}

	cblogger.Info(result)
	spew.Dump(result)

	if len(result.RouteTables) > 0 {
		routeTableId := *result.RouteTables[0].RouteTableId
		cblogger.Infof("라우팅 테이블 ID 찾음 : [%s]", routeTableId)
		return routeTableId, nil
	} else {
		return "", errors.New("VPC에 할당된 라우팅 테이블 ID를 찾을 수 없습니다.")
	}
}

func (VPCHandler *TencentVPCHandler) CreateSubnet(vpcId string, reqSubnetInfo irs.SubnetInfo) (irs.SubnetInfo, error) {
	cblogger.Info(reqSubnetInfo)

	zoneId := VPCHandler.Region.Zone
	cblogger.Infof("Zone : %s", zoneId)
	if zoneId == "" {
		cblogger.Error("Connection 정보에 Zone 정보가 없습니다.")
		return irs.SubnetInfo{}, errors.New("Connection 정보에 Zone 정보가 없습니다.")
	}

	if reqSubnetInfo.IId.SystemId != "" {
		vpcInfo, errVpcInfo := VPCHandler.GetSubnet(reqSubnetInfo.IId.SystemId)
		if errVpcInfo == nil {
			cblogger.Errorf("이미 [%S] Subnet이 존재하기 때문에 생성하지 않고 기존 정보와 함께 에러를 리턴함.", reqSubnetInfo.IId.SystemId)
			cblogger.Info(vpcInfo)
			return vpcInfo, errors.New("InvalidVNetwork.Duplicate: The Subnet '" + reqSubnetInfo.IId.SystemId + "' already exists.")
		}
	}

	//서브넷 생성
	input := &ec2.CreateSubnetInput{
		CidrBlock: aws.String(reqSubnetInfo.IPv4_CIDR),
		VpcId:     aws.String(vpcId),
		//AvailabilityZoneId: aws.String(zoneId),	//use1-az1, use1-az2, use1-az3, use1-az4, use1-az5, use1-az6
		AvailabilityZone: aws.String(zoneId),
	}

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")

	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.AWS,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: reqSubnetInfo.IId.NameId,
		CloudOSAPI:   "CreateSubnet()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	start := call.Start()

	cblogger.Info(input)
	result, err := VPCHandler.Client.CreateSubnet(input)
	callLogInfo.ElapsedTime = call.Elapsed(start)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
				callLogInfo.ErrorMSG = aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
			callLogInfo.ErrorMSG = err.Error()
		}
		callogger.Info(call.String(callLogInfo))
		return irs.SubnetInfo{}, err
	}
	callogger.Info(call.String(callLogInfo))
	cblogger.Info(result)
	//spew.Dump(result)

	//vNetworkInfo := irs.VNetworkInfo{}
	vNetworkInfo := ExtractSubnetDescribeInfo(result.Subnet)

	//Subnet Name 태깅
	if SetNameTag(VPCHandler.Client, *result.Subnet.SubnetId, reqSubnetInfo.IId.NameId) {
		cblogger.Infof("Subnet에 %s Name 설정 성공", reqSubnetInfo.IId.NameId)
	} else {
		cblogger.Errorf("Subnet에 %s Name 설정 실패", reqSubnetInfo.IId.NameId)
	}

	vNetworkInfo.IId.NameId = reqSubnetInfo.IId.NameId

	// VPC의 라우팅 테이블에 생성된 Subnet 정보를 추가 함.
	errSubnetRoute := VPCHandler.AssociateRouteTable(vpcId, vNetworkInfo.IId.SystemId)
	if errSubnetRoute != nil {
	} else {
		return vNetworkInfo, errSubnetRoute
	}

	return vNetworkInfo, nil
}

// VPC의 라우팅 테이블에 생성된 Subnet을 연결 함.
func (VPCHandler *TencentVPCHandler) AssociateRouteTable(vpcId string, subnetId string) error {
	routeTableId, errRoute := VPCHandler.GetDefaultRouteTable(vpcId)
	if errRoute != nil {
		return errRoute
	}

	input := &ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(routeTableId),
		SubnetId:     aws.String(subnetId),
	}

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")

	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.AWS,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: subnetId,
		CloudOSAPI:   "AssociateRouteTable()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	callLogStart := call.Start()

	result, err := VPCHandler.Client.AssociateRouteTable(input)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
				callLogInfo.ErrorMSG = aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
			callLogInfo.ErrorMSG = err.Error()
		}
		callogger.Info(call.String(callLogInfo))
		return err
	}

	callogger.Info(call.String(callLogInfo))
	cblogger.Info(result)
	//spew.Dump(result)
	return nil
}

func (VPCHandler *TencentVPCHandler) ListVPC() ([]*irs.VPCInfo, error) {
	cblogger.Debug("Start")
	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.AWS,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: "ListVPC",
		CloudOSAPI:   "DescribeVpcs()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	callLogStart := call.Start()

	result, err := VPCHandler.Client.DescribeVpcs(&ec2.DescribeVpcsInput{})
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
		}
		callLogInfo.ErrorMSG = err.Error()
		callogger.Info(call.String(callLogInfo))
		return nil, err
	}
	callogger.Info(call.String(callLogInfo))

	var vNetworkInfoList []*irs.VPCInfo
	for _, curVpc := range result.Vpcs {
		cblogger.Infof("[%s] VPC 정보 조회", *curVpc.VpcId)
		vNetworkInfo, vpcErr := VPCHandler.GetVPC(irs.IID{SystemId: *curVpc.VpcId})
		if vpcErr != nil {
			return nil, vpcErr
		}
		vNetworkInfoList = append(vNetworkInfoList, &vNetworkInfo)
	}

	spew.Dump(vNetworkInfoList)
	return vNetworkInfoList, nil
}

func (VPCHandler *TencentVPCHandler) GetVPC(vpcIID irs.IID) (irs.VPCInfo, error) {
	cblogger.Info("VPC IID : ", vpcIID.SystemId)

	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{
			aws.String(vpcIID.SystemId),
		},
	}

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.AWS,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: vpcIID.SystemId,
		CloudOSAPI:   "DescribeVpcs()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	callLogStart := call.Start()

	result, err := VPCHandler.Client.DescribeVpcs(input)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
				callLogInfo.ErrorMSG = aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
			callLogInfo.ErrorMSG = err.Error()
		}
		callogger.Info(call.String(callLogInfo))
		return irs.VPCInfo{}, err
	}
	callogger.Info(call.String(callLogInfo))

	cblogger.Info(result)
	//spew.Dump(result)

	if reflect.ValueOf(result.Vpcs).IsNil() {
		return irs.VPCInfo{}, nil
	}

	var errSubnet error
	awsVpcInfo := ExtractVpcDescribeInfo(result.Vpcs[0])
	awsVpcInfo.SubnetInfoList, errSubnet = VPCHandler.ListSubnet(vpcIID.SystemId)
	if errSubnet != nil {
		return awsVpcInfo, errSubnet
	}

	return awsVpcInfo, nil
}

//VPC 정보를 추출함
func ExtractVpcDescribeInfo(vpcInfo *ec2.Vpc) irs.VPCInfo {
	awsVpcInfo := irs.VPCInfo{
		IId:       irs.IID{SystemId: *vpcInfo.VpcId},
		IPv4_CIDR: *vpcInfo.CidrBlock,
		//IsDefault: *vpcInfo.IsDefault,
		//State:     *vpcInfo.State,
	}

	//Name은 Tag의 "Name" 속성에만 저장됨
	//NameId는 전달할 필요가 없음.
	return awsVpcInfo
}

func (VPCHandler *TencentVPCHandler) DeleteSubnet(subnetIID irs.IID) (bool, error) {
	input := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(subnetIID.SystemId),
	}
	cblogger.Info(input)

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")

	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.AWS,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: subnetIID.SystemId,
		CloudOSAPI:   "DeleteSubnet()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	start := call.Start()

	_, err := VPCHandler.Client.DeleteSubnet(input)
	callLogInfo.ElapsedTime = call.Elapsed(start)
	cblogger.Info(err)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
				callLogInfo.ErrorMSG = aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and Message from an error.
			cblogger.Error(err.Error())
			callLogInfo.ErrorMSG = err.Error()
		}
		callogger.Info(call.String(callLogInfo))
		return false, err
	}

	callogger.Info(call.String(callLogInfo))
	return true, nil
}

func (VPCHandler *TencentVPCHandler) DeleteVPC(vpcIID irs.IID) (bool, error) {
	cblogger.Infof("Delete VPC : [%s]", vpcIID.SystemId)

	vpcInfo, errVpcInfo := VPCHandler.GetVPC(vpcIID)
	if errVpcInfo != nil {
		return false, errVpcInfo
	}

	//=================
	// Subnet삭제
	//=================
	for _, curSubnet := range vpcInfo.SubnetInfoList {
		cblogger.Infof("[%s] Subnet 삭제", curSubnet.IId.SystemId)
		delSubnet, errSubnet := VPCHandler.DeleteSubnet(curSubnet.IId)
		if errSubnet != nil {
			return false, errSubnet
		}

		if delSubnet {
			cblogger.Infof("  ==> [%s] Subnet 삭제완료", curSubnet.IId.SystemId)
		} else {
			cblogger.Errorf("  ==> [%s] Subnet 삭제실패", curSubnet.IId.SystemId)
			return false, errors.New("Subnet 삭제 실패로 VPC를 삭제하지 못 했습니다.") //삭제 실패 이유를 모르는 경우
		}
	}

	cblogger.Infof("[%s] VPC를 삭제 함.", vpcInfo.IId.SystemId)
	cblogger.Info("VPC 제거를 위해 생성된 IGW / Route들 제거 시작")

	// 라우팅 테이블에 추가한 IGW 라우터를 먼저 삭제함.
	errRoute := VPCHandler.DeleteRouteIGW(vpcInfo.IId.SystemId)
	if errRoute != nil {
		cblogger.Error("라우팅 테이블에 추가한 0.0.0.0/0 IGW 라우터 삭제 실패")
		cblogger.Error(errRoute)
		if "InvalidRoute.NotFound" == errRoute.Error() {
			cblogger.Infof("[%s]예외는 #255예외에 의해 정상으로 간주하고 다음 단계를 진행함.", errRoute)
		} else {
			return false, errRoute
		}
		//} else {
		//	cblogger.Info("라우팅 테이블에 추가한 0.0.0.0/0 IGW 라우터 삭제 완료")
	}

	//VPC에 연결된 모든 IGW를 삭제함. (VPC에 할당된 모든 IGW조회후 삭제)
	errIgw := VPCHandler.DeleteAllIGW(vpcInfo.IId.SystemId)
	if errIgw != nil {
		cblogger.Error("모든 IGW 삭제 실패 : ", errIgw)
	} else {
		cblogger.Info("모든 IGW 삭제 완료")
	}

	input := &ec2.DeleteVpcInput{
		VpcId: aws.String(vpcInfo.IId.SystemId),
	}

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.AWS,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: vpcInfo.IId.SystemId,
		CloudOSAPI:   "DeleteVpc()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	callLogStart := call.Start()

	result, err := VPCHandler.Client.DeleteVpc(input)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
				callLogInfo.ErrorMSG = aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and Message from an error.
			cblogger.Error(err.Error())
			callLogInfo.ErrorMSG = err.Error()
		}
		callogger.Info(call.String(callLogInfo))
		return false, err
	}
	callogger.Info(call.String(callLogInfo))

	cblogger.Info(result)
	spew.Dump(result)
	return true, nil
}

// VPC에 설정된 0.0.0.0/0 라우터를 제거 함.
// #255예외 처리 보완에 따른 라우팅 정보 삭제전 0.0.0.0 조회후 삭제하도록 로직 변경
func (VPCHandler *TencentVPCHandler) DeleteRouteIGW(vpcId string) error {
	cblogger.Infof("VPC ID : [%s]", vpcId)
	routeTableId := ""

	input := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(vpcId),
				},
			},
		},
	}

	result, err := VPCHandler.Client.DescribeRouteTables(input)
	if err != nil {
		return err
	}

	cblogger.Info(result)
	spew.Dump(result)

	if len(result.RouteTables) < 1 {
		return errors.New("VPC에 할당된 라우팅 테이블 정보를 찾을 수 없습니다.")
	}

	routeTableId = *result.RouteTables[0].RouteTableId
	cblogger.Infof("라우팅 테이블 ID 찾음 : [%s]", routeTableId)

	cblogger.Infof("RouteTable[%s]에 할당된 라우팅(0.0.0.0/0) 정보를 조회합니다.", routeTableId)

	//ec2.Route
	findIgw := false
	for _, curRoute := range result.RouteTables[0].Routes {
		cblogger.Infof("DestinationCidrBlock[%s] Check", *curRoute.DestinationCidrBlock)

		if "0.0.0.0/0" == *curRoute.DestinationCidrBlock {
			cblogger.Infof("===>RouteTable[%s]에 할당된 라우팅(0.0.0.0/0) 정보를 찾았습니다!!", routeTableId)
			findIgw = true
			break
		}
	}

	if !findIgw {
		cblogger.Infof("RouteTable[%s]에 할당된 IGW의 라우팅(0.0.0.0/0) 정보가 없으므로 라우트 삭제처리는 중단합니다. ", routeTableId)
		return nil
	}

	cblogger.Infof("RouteTable[%s]에 할당된 라우팅(0.0.0.0/0) 정보를 삭제합니다.", routeTableId)
	inputDel := &ec2.DeleteRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		RouteTableId:         aws.String(routeTableId),
	}
	cblogger.Info(inputDel)

	//https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DeleteRoute.html
	resultDel, err := VPCHandler.Client.DeleteRoute(inputDel)
	if err != nil {
		cblogger.Errorf("RouteTable[%s]에 대한 라우팅(0.0.0.0/0) 정보 삭제 실패", routeTableId)
		if aerr, ok := err.(awserr.Error); ok {
			//InvalidRoute.NotFound
			cblogger.Errorf("Error Code : [%s] - Error:[%s] - Message:[%s]", aerr.Code(), aerr.Error(), aerr.Message())
			switch aerr.Code() {
			case "InvalidRoute.NotFound": //NotFound에러는 무시하라고 해서 (예외#255)
				return errors.New(aerr.Code())
			default:
				cblogger.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
		}
		return err
	}
	cblogger.Infof("RouteTable[%s]에 대한 라우팅(0.0.0.0/0) 정보 삭제 완료", routeTableId)

	cblogger.Info(resultDel)
	spew.Dump(resultDel)
	cblogger.Info("라우팅 테이블에 추가한 0.0.0.0/0 IGW 라우터 삭제 완료")
	return nil
}

//VPC에 연결된 모든 IGW를 삭제함.
func (VPCHandler *TencentVPCHandler) DeleteAllIGW(vpcId string) error {
	input := &ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("attachment.vpc-id"),
				Values: []*string{
					aws.String(vpcId),
				},
			},
		},
	}

	result, err := VPCHandler.Client.DescribeInternetGateways(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
		}
		return err
	}

	cblogger.Info(result)
	spew.Dump(result)

	// VPC 삭제를 위해 연결된 모든 IGW 제거
	// 일단, 에러는 무시함.
	for _, curIgw := range result.InternetGateways {
		//IGW 삭제전 연결된 IGW의 연결을 끊어야함.
		VPCHandler.DetachInternetGateway(vpcId, *curIgw.InternetGatewayId)
		//IGW 삭제
		VPCHandler.DeleteIGW(*curIgw.InternetGatewayId)
	}

	return nil
}

// VPC에 연결된 IGW의 연결을 해제함.
func (VPCHandler *TencentVPCHandler) DetachInternetGateway(vpcId string, igwId string) error {
	cblogger.Infof("VPC[%s]에 연결된 IGW[%s]의 연결을 해제함.", vpcId, igwId)

	input := &ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String(igwId),
		VpcId:             aws.String(vpcId),
	}

	result, err := VPCHandler.Client.DetachInternetGateway(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
		}
		return err
	}

	cblogger.Info(result)
	spew.Dump(result)
	return nil
}

// IGW를 삭제 함.
func (VPCHandler *TencentVPCHandler) DeleteIGW(igwId string) error {
	input := &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(igwId),
	}

	result, err := VPCHandler.Client.DeleteInternetGateway(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
		}
		return err
	}

	cblogger.Info(result)
	spew.Dump(result)
	return nil
}

//VPC의 하위 서브넷 목록을 조회함.
func (VPCHandler *TencentVPCHandler) ListSubnet(vpcId string) ([]irs.SubnetInfo, error) {
	cblogger.Debug("Start")
	var arrSubnetInfoList []irs.SubnetInfo

	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(vpcId),
				},
			},
		},
	}

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")

	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.AWS,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: "ListSubnet",
		CloudOSAPI:   "DescribeSubnets()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	callLogStart := call.Start()

	//spew.Dump(input)
	result, err := VPCHandler.Client.DescribeSubnets(input)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
				callLogInfo.ErrorMSG = aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
			callLogInfo.ErrorMSG = err.Error()
		}
		callogger.Info(call.String(callLogInfo))
		return nil, err
	}
	callogger.Info(call.String(callLogInfo))

	spew.Dump(result)
	for _, curSubnet := range result.Subnets {
		cblogger.Infof("[%s] Subnet 정보 조회", *curSubnet.SubnetId)
		arrSubnetInfo := ExtractSubnetDescribeInfo(curSubnet)
		arrSubnetInfoList = append(arrSubnetInfoList, arrSubnetInfo)
	}

	spew.Dump(arrSubnetInfoList)
	return arrSubnetInfoList, nil
}

func (VPCHandler *TencentVPCHandler) GetSubnet(reqSubnetId string) (irs.SubnetInfo, error) {
	cblogger.Infof("SubnetId : [%s]", reqSubnetId)

	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []*string{
			aws.String(reqSubnetId),
		},
	}

	spew.Dump(input)
	// logger for HisCall
	callogger := call.GetLogger("HISCALL")

	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.AWS,
		RegionZone:   VPCHandler.Region.Zone,
		ResourceType: call.VPCSUBNET,
		ResourceName: reqSubnetId,
		CloudOSAPI:   "DescribeSubnets()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	callLogStart := call.Start()
	result, err := VPCHandler.Client.DescribeSubnets(input)
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)
	cblogger.Info(result)
	//spew.Dump(result)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				cblogger.Error(aerr.Error())
				callLogInfo.ErrorMSG = aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			cblogger.Error(err.Error())
			callLogInfo.ErrorMSG = err.Error()
		}
		callogger.Info(call.String(callLogInfo))
		return irs.SubnetInfo{}, err
	}
	callogger.Info(call.String(callLogInfo))

	if !reflect.ValueOf(result.Subnets).IsNil() {
		retSubnetInfo := ExtractSubnetDescribeInfo(result.Subnets[0])
		return retSubnetInfo, nil
	} else {
		return irs.SubnetInfo{}, errors.New("InvalidSubnet.NotFound: The CBVnetwork '" + reqSubnetId + "' does not exist")
	}
}


//Subnet 정보를 추출함
func ExtractSubnetDescribeInfo(subnetInfo *ec2.Subnet) irs.SubnetInfo {
	vNetworkInfo := irs.SubnetInfo{
		IId:       irs.IID{SystemId: *subnetInfo.SubnetId},
		IPv4_CIDR: *subnetInfo.CidrBlock,
		//Status:    *subnetInfo.State,
	}

	keyValueList := []irs.KeyValue{
		{Key: "VpcId", Value: *subnetInfo.VpcId},
		{Key: "MapPublicIpOnLaunch", Value: strconv.FormatBool(*subnetInfo.MapPublicIpOnLaunch)},
		{Key: "AvailableIpAddressCount", Value: strconv.FormatInt(*subnetInfo.AvailableIpAddressCount, 10)},
		{Key: "AvailabilityZone", Value: *subnetInfo.AvailabilityZone},
		{Key: "Status", Value: *subnetInfo.State},
	}
	vNetworkInfo.KeyValueList = keyValueList

	return vNetworkInfo
}

func (VPCHandler *TencentVPCHandler) AddSubnet(vpcIID irs.IID, subnetInfo irs.SubnetInfo) (irs.VPCInfo, error) {
	cblogger.Infof("[%s] Subnet 추가 - CIDR : %s", subnetInfo.IId.NameId, subnetInfo.IPv4_CIDR)
	resSubnet, errSubnet := VPCHandler.CreateSubnet(vpcIID.SystemId, subnetInfo)
	if errSubnet != nil {
		cblogger.Error(errSubnet)
		return irs.VPCInfo{}, errSubnet
	}
	cblogger.Info(resSubnet)

	vpcInfo, errVpcInfo := VPCHandler.GetVPC(vpcIID)
	if errVpcInfo != nil {
		cblogger.Error(errVpcInfo)
		return irs.VPCInfo{}, errVpcInfo
	}

	findSubnet := false
	cblogger.Debug("============== 체크할 값 =========")
	for posSubnet, curSubnetInfo := range vpcInfo.SubnetInfoList {
		cblogger.Debugf("%d - [%s] Subnet 처리 시작", posSubnet, curSubnetInfo.IId.SystemId)
		if resSubnet.IId.SystemId == curSubnetInfo.IId.SystemId {
			cblogger.Infof("추가 요청 받은 [%s] Subnet을 발견 했습니다. - SystemID:[%s]", subnetInfo.IId.NameId, curSubnetInfo.IId.SystemId)
			//for ~ range는 포인터가 아니라서 값 수정이 안됨. for loop으로 직접 서브넷을 체크하거나 vpcInfo의 배열의 값을 수정해야 함.
			cblogger.Infof("인덱스 위치 : %d", posSubnet)
			//vpcInfo.SubnetInfoList[posSubnet].IId.NameId = "테스트~"
			vpcInfo.SubnetInfoList[posSubnet].IId.NameId = subnetInfo.IId.NameId
			findSubnet = true
			break
		}
	}

	if !findSubnet {
		cblogger.Errorf("서브넷 생성은 성공했으나 VPC의 서브넷 목록에서 추가 요청한 [%s]서브넷의 정보[%s]를 찾지 못했습니다.", subnetInfo.IId.NameId, resSubnet.IId.SystemId)
		return irs.VPCInfo{}, errors.New("MismatchSubnet.NotFound: No SysmteId[" + resSubnet.IId.SystemId + "] found for newly created Subnet[" + subnetInfo.IId.NameId + "].")
	}
	//spew.Dump(vpcInfo)

	return vpcInfo, nil
}

func (VPCHandler *TencentVPCHandler) RemoveSubnet(vpcIID irs.IID, subnetIID irs.IID) (bool, error) {
	cblogger.Infof("[%s] VPC의 [%s] Subnet 삭제", vpcIID.SystemId, subnetIID.SystemId)

	return VPCHandler.DeleteSubnet(subnetIID)
	//return false, nil
}
*/
