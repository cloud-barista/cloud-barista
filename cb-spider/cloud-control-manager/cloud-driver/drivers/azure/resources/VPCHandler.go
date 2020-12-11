package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/Azure/go-autorest/autorest/to"

	call "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/call-log"
	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
)

const (
	VPC = "VPC"
)

type AzureVPCHandler struct {
	Region       idrv.RegionInfo
	Ctx          context.Context
	Client       *network.VirtualNetworksClient
	SubnetClient *network.SubnetsClient
}

func (vpcHandler *AzureVPCHandler) setterVPC(network network.VirtualNetwork) *irs.VPCInfo {
	vpcInfo := &irs.VPCInfo{
		IId: irs.IID{
			NameId:   *network.Name,
			SystemId: *network.ID,
		},
		IPv4_CIDR:    (*network.AddressSpace.AddressPrefixes)[0],
		KeyValueList: []irs.KeyValue{{Key: "ResourceGroup", Value: vpcHandler.Region.ResourceGroup}},
	}

	subnetArr := make([]irs.SubnetInfo, len(*network.Subnets))
	for i, subnet := range *network.Subnets {
		subnetArr[i] = *vpcHandler.setterSubnet(subnet)
	}
	vpcInfo.SubnetInfoList = subnetArr

	return vpcInfo
}

func (vpcHandler *AzureVPCHandler) setterSubnet(subnet network.Subnet) *irs.SubnetInfo {
	subnetInfo := &irs.SubnetInfo{
		IId: irs.IID{
			NameId:   *subnet.Name,
			SystemId: *subnet.ID,
		},
		IPv4_CIDR:    *subnet.AddressPrefix,
		KeyValueList: []irs.KeyValue{{Key: "ResourceGroup", Value: vpcHandler.Region.ResourceGroup}},
	}
	return subnetInfo
}

func (vpcHandler *AzureVPCHandler) CreateVPC(vpcReqInfo irs.VPCReqInfo) (irs.VPCInfo, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vpcHandler.Region, call.VPCSUBNET, VPC, "CreateVPC()")

	// Check VPC Exists
	vpc, _ := vpcHandler.Client.Get(vpcHandler.Ctx, vpcHandler.Region.ResourceGroup, vpcReqInfo.IId.NameId, "")
	if vpc.ID != nil {
		createErr := errors.New(fmt.Sprintf("vpc with name %s already exist", vpcReqInfo.IId.NameId))
		LoggingError(hiscallInfo, createErr)
		return irs.VPCInfo{}, createErr
	}

	// Create VPC
	createOpts := network.VirtualNetwork{
		Name: to.StringPtr(vpcReqInfo.IId.NameId),
		VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
			AddressSpace: &network.AddressSpace{
				AddressPrefixes: &[]string{vpcReqInfo.IPv4_CIDR},
			},
		},
		Location: &vpcHandler.Region.Region,
	}

	start := call.Start()
	future, err := vpcHandler.Client.CreateOrUpdate(vpcHandler.Ctx, vpcHandler.Region.ResourceGroup, vpcReqInfo.IId.NameId, createOpts)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.VPCInfo{}, err
	}
	err = future.WaitForCompletionRef(vpcHandler.Ctx, vpcHandler.Client.Client)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.VPCInfo{}, err
	}
	LoggingInfo(hiscallInfo, start)

	// Create Subnet
	var subnetCreateOpts network.Subnet
	for _, subnet := range vpcReqInfo.SubnetInfoList {
		subnetCreateOpts = network.Subnet{
			Name: to.StringPtr(subnet.IId.NameId),
			SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr(subnet.IPv4_CIDR),
			},
		}
		future, err := vpcHandler.SubnetClient.CreateOrUpdate(vpcHandler.Ctx, vpcHandler.Region.ResourceGroup, vpcReqInfo.IId.NameId, subnet.IId.NameId, subnetCreateOpts)
		if err != nil {
			cblogger.Error(fmt.Sprintf("failed to create subnet with name %s", subnet.IId.NameId))
			continue
		}
		err = future.WaitForCompletionRef(vpcHandler.Ctx, vpcHandler.Client.Client)
		if err != nil {
			cblogger.Error(fmt.Sprintf("failed to create subnet with name %s", subnet.IId.NameId))
			continue
		}
	}

	// 생성된 VNetwork 정보 리턴
	vpcInfo, err := vpcHandler.GetVPC(irs.IID{NameId: vpcReqInfo.IId.NameId})
	if err != nil {
		return irs.VPCInfo{}, err
	}
	return vpcInfo, nil
}

func (vpcHandler *AzureVPCHandler) ListVPC() ([]*irs.VPCInfo, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vpcHandler.Region, call.VPCSUBNET, VPC, "ListVPC()")

	start := call.Start()
	vpcList, err := vpcHandler.Client.List(vpcHandler.Ctx, vpcHandler.Region.ResourceGroup)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return nil, err
	}
	LoggingInfo(hiscallInfo, start)

	vpcInfoList := make([]*irs.VPCInfo, len(vpcList.Values()))
	for i, vpc := range vpcList.Values() {
		vpcInfoList[i] = vpcHandler.setterVPC(vpc)
	}
	return vpcInfoList, nil
}

func (vpcHandler *AzureVPCHandler) GetVPC(vpcIID irs.IID) (irs.VPCInfo, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vpcHandler.Region, call.VPCSUBNET, vpcIID.NameId, "GetVPC()")

	start := call.Start()
	vpc, err := vpcHandler.Client.Get(vpcHandler.Ctx, vpcHandler.Region.ResourceGroup, vpcIID.NameId, "")
	if err != nil {
		LoggingError(hiscallInfo, err)
		return irs.VPCInfo{}, err
	}
	LoggingInfo(hiscallInfo, start)

	vpcInfo := vpcHandler.setterVPC(vpc)
	return *vpcInfo, nil
}

func (vpcHandler *AzureVPCHandler) DeleteVPC(vpcIID irs.IID) (bool, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vpcHandler.Region, call.VPCSUBNET, vpcIID.NameId, "DeleteVPC()")

	start := call.Start()
	future, err := vpcHandler.Client.Delete(vpcHandler.Ctx, vpcHandler.Region.ResourceGroup, vpcIID.NameId)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return false, err
	}
	err = future.WaitForCompletionRef(vpcHandler.Ctx, vpcHandler.Client.Client)
	if err != nil {
		LoggingError(hiscallInfo, err)
		return false, err
	}
	LoggingInfo(hiscallInfo, start)

	return true, nil
}

func (vpcHandler *AzureVPCHandler) AddSubnet(vpcIID irs.IID, subnetInfo irs.SubnetInfo) (irs.VPCInfo, error) {
	// log HisCall
	hiscallInfo := GetCallLogScheme(vpcHandler.Region, call.VPCSUBNET, subnetInfo.IId.NameId, "AddSubnet()")

	vpc, _ := vpcHandler.Client.Get(vpcHandler.Ctx, vpcHandler.Region.ResourceGroup, vpcIID.NameId, "")
	if vpc.ID == nil {
		createErr := errors.New(fmt.Sprintf("vpc with name %s not exist", vpcIID.NameId))
		LoggingError(hiscallInfo, createErr)
		return irs.VPCInfo{}, createErr
	}
	subnetCreateOpts := network.Subnet{
		Name: to.StringPtr(subnetInfo.IId.NameId),
		SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr(subnetInfo.IPv4_CIDR),
		},
	}
	future, err := vpcHandler.SubnetClient.CreateOrUpdate(vpcHandler.Ctx, vpcHandler.Region.ResourceGroup, *vpc.Name, subnetInfo.IId.NameId, subnetCreateOpts)
	if err != nil {
		createErr := errors.New(fmt.Sprintf("failed to create subnet with name %s", subnetInfo.IId.NameId))
		LoggingError(hiscallInfo, createErr)
		return irs.VPCInfo{}, createErr
	}
	err = future.WaitForCompletionRef(vpcHandler.Ctx, vpcHandler.Client.Client)
	if err != nil {
		createErr := errors.New(fmt.Sprintf("failed to create subnet with name %s", subnetInfo.IId.NameId))
		LoggingError(hiscallInfo, createErr)
		return irs.VPCInfo{}, createErr
	}
	result, err := vpcHandler.GetVPC(irs.IID{NameId: vpcIID.NameId})
	if err != nil {
		createErr := errors.New(fmt.Sprintf("failed to create subnet with name %s", subnetInfo.IId.NameId))
		LoggingError(hiscallInfo, createErr)
		return irs.VPCInfo{}, createErr
	}
	return result, nil
}

func (vpcHandler *AzureVPCHandler) RemoveSubnet(vpcIID irs.IID, subnetIID irs.IID) (bool, error) {
	future, err := vpcHandler.SubnetClient.Delete(vpcHandler.Ctx, vpcHandler.Region.ResourceGroup, vpcIID.NameId, subnetIID.NameId)
	if err != nil {
		return false, err
	}
	err = future.WaitForCompletionRef(vpcHandler.Ctx, vpcHandler.Client.Client)
	if err != nil {
		return false, err
	}
	return true, nil
}
