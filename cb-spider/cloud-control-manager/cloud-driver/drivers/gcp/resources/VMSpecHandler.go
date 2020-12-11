// Proof of Concepts of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is a Cloud Driver Example for PoC Test.
//
// program by ysjeon@mz.co.kr, 2020.01.

package resources

import (
	"context"
	_ "errors"
	"strconv"
	"strings"

	compute "google.golang.org/api/compute/v1"

	call "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/call-log"
	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
)

type GCPVMSpecHandler struct {
	Region     idrv.RegionInfo
	Ctx        context.Context
	Client     *compute.Service
	Credential idrv.CredentialInfo
}

func (vmSpecHandler *GCPVMSpecHandler) ListVMSpec(Region string) ([]*irs.VMSpecInfo, error) {

	projectID := vmSpecHandler.Credential.ProjectID
	zone := vmSpecHandler.Region.Zone

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.GCP,
		RegionZone:   vmSpecHandler.Region.Zone,
		ResourceType: call.VMSPEC,
		ResourceName: "",
		CloudOSAPI:   "List()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	callLogStart := call.Start()
	resp, err := vmSpecHandler.Client.MachineTypes.List(projectID, zone).Do()
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	if err != nil {
		callLogInfo.ErrorMSG = err.Error()
		callogger.Info(call.String(callLogInfo))
		cblogger.Error(err)
		return []*irs.VMSpecInfo{}, err
	}
	callogger.Info(call.String(callLogInfo))
	var vmSpecInfo []*irs.VMSpecInfo
	for _, i := range resp.Items {
		info := irs.VMSpecInfo{
			Region: zone,
			Name:   i.Name,
			VCpu: irs.VCpuInfo{
				Count: strconv.FormatInt(i.GuestCpus, 10),
			},
			Mem: strconv.FormatInt(i.MemoryMb, 10),
		}
		vmSpecInfo = append(vmSpecInfo, &info)
	}
	return vmSpecInfo, nil
}

func (vmSpecHandler *GCPVMSpecHandler) GetVMSpec(Region string, Name string) (irs.VMSpecInfo, error) {
	// default info
	projectID := vmSpecHandler.Credential.ProjectID
	zone := vmSpecHandler.Region.Zone

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.GCP,
		RegionZone:   vmSpecHandler.Region.Zone,
		ResourceType: call.VMSPEC,
		ResourceName: Name,
		CloudOSAPI:   "Get()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	callLogStart := call.Start()

	info, err := vmSpecHandler.Client.MachineTypes.Get(projectID, zone, Name).Do()
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	if err != nil {
		callLogInfo.ErrorMSG = err.Error()
		callogger.Info(call.String(callLogInfo))
		cblogger.Error(err)
		return irs.VMSpecInfo{}, err
	}
	callogger.Info(call.String(callLogInfo))

	vmSpecInfo := irs.VMSpecInfo{
		Region: Region,
		Name:   Name,
		VCpu: irs.VCpuInfo{
			Count: strconv.FormatInt(info.GuestCpus, 10),
			Clock: "",
		},
		Mem: strconv.FormatInt(info.MemoryMb, 10),
		Gpu: []irs.GpuInfo{
			{
				Count: "",
				Mfr:   "",
				Model: "",
				Mem:   "",
			},
		},
	}

	return vmSpecInfo, nil
}

func (vmSpecHandler *GCPVMSpecHandler) ListOrgVMSpec(Region string) (string, error) {
	projectID := vmSpecHandler.Credential.ProjectID
	zone := vmSpecHandler.Region.Zone

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.GCP,
		RegionZone:   vmSpecHandler.Region.Zone,
		ResourceType: call.VMSPEC,
		ResourceName: "",
		CloudOSAPI:   "List()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	callLogStart := call.Start()
	resp, err := vmSpecHandler.Client.MachineTypes.List(projectID, zone).Do()
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	if err != nil {
		callLogInfo.ErrorMSG = err.Error()
		callogger.Info(call.String(callLogInfo))
		cblogger.Error(err)
		return "", err
	}
	callogger.Info(call.String(callLogInfo))
	j, _ := resp.MarshalJSON()

	return string(j), err
}

func (vmSpecHandler *GCPVMSpecHandler) GetOrgVMSpec(Region string, Name string) (string, error) {
	projectID := vmSpecHandler.Credential.ProjectID
	zone := vmSpecHandler.Region.Zone

	// logger for HisCall
	callogger := call.GetLogger("HISCALL")
	callLogInfo := call.CLOUDLOGSCHEMA{
		CloudOS:      call.GCP,
		RegionZone:   vmSpecHandler.Region.Zone,
		ResourceType: call.VMSPEC,
		ResourceName: Name,
		CloudOSAPI:   "Get()",
		ElapsedTime:  "",
		ErrorMSG:     "",
	}
	callLogStart := call.Start()

	info, err := vmSpecHandler.Client.MachineTypes.Get(projectID, zone, Name).Do()
	callLogInfo.ElapsedTime = call.Elapsed(callLogStart)

	if err != nil {
		callLogInfo.ErrorMSG = err.Error()
		callogger.Info(call.String(callLogInfo))
		cblogger.Error(err)
		return "", err
	}
	callogger.Info(call.String(callLogInfo))
	j, _ := info.MarshalJSON()

	return string(j), err
}

// gcp 같은경우 n1 타입만 그래픽 카드가 추가 되며
// 1. n1타입인지 확인하는 로직 필요
// 2. 해당 카드에 관련된 정보를 조회하는 로직필요.
// 3. 해당 리스트를 조회하고 해당 GPU를 선택하는 로직

func CheckMachineType(Name string) bool {
	prefix := "n1"

	if ok := strings.HasPrefix(prefix, Name); ok {
		return ok
	}

	return false

}
