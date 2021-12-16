// Cloud Driver Interface of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is Mock Driver.
//
// by CB-Spider Team, 2020.09.

package resources

import (
	"fmt"

	cblog "github.com/cloud-barista/cb-log"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
)

var vmSpecInfoMap map[string][]*irs.VMSpecInfo

type MockVMSpecHandler struct {
	MockName string
}

var prepareVMSpecInfoList []*irs.VMSpecInfo

func init() {
	vmSpecInfoMap = make(map[string][]*irs.VMSpecInfo)
}

// Be called before using the User function.
// Called in MockDriver
func PrepareVMSpec(mockName string) {
        cblogger := cblog.GetLogger("CB-SPIDER")
        cblogger.Info("Mock Driver: called prepare()!")

        if vmSpecInfoMap[mockName] != nil {
        	return
		}

		prepareVMSpecInfoList = []*irs.VMSpecInfo{
                {"common-region", "mock-vmspec-01", irs.VCpuInfo{"4", "2.7"}, "32768", []irs.GpuInfo{{"2", "NVIDIA", "V100", "16384MB"}}, nil},
                {"common-region", "mock-vmspec-02", irs.VCpuInfo{"4", "3.2"}, "32768", []irs.GpuInfo{{"1", "NVIDIA", "V100", "16384MB"}}, nil},
                {"common-region", "mock-vmspec-03", irs.VCpuInfo{"8", "2.7"}, "62464", nil, nil},
                {"common-region", "mock-vmspec-04", irs.VCpuInfo{"8", "2.7"}, "1024", nil, nil},
        }
		vmSpecInfoMap[mockName] = prepareVMSpecInfoList
}

func (vmSpecHandler *MockVMSpecHandler) ListVMSpec() ([]*irs.VMSpecInfo, error) {
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called ListVMSpec()!")

	mockName := vmSpecHandler.MockName

	infoList, ok := vmSpecInfoMap[mockName]
	if !ok {
		return []*irs.VMSpecInfo{}, nil
	}

	// cloning list of VMSpec
	resultList := make([]*irs.VMSpecInfo, len(infoList))
	copy(resultList, infoList)
	return resultList, nil
}

func (vmSpecHandler *MockVMSpecHandler) GetVMSpec(Name string) (irs.VMSpecInfo, error) {
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called GetVMSpec()!")

	infoList, err := vmSpecHandler.ListVMSpec()
	if err != nil {
		cblogger.Error(err)
		return irs.VMSpecInfo{}, err
	}

	for _, info := range infoList {
		if (*info).Name == Name {
			return *info, nil
		}
	}

	return irs.VMSpecInfo{}, fmt.Errorf("%s VMSpec does not exist!!", Name)
}

func (vmSpecHandler *MockVMSpecHandler) ListOrgVMSpec() (string, error) { // return string: json format
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called ListOrgVMSpec()!")
	return "", nil
}

func (vmSpecHandler *MockVMSpecHandler) GetOrgVMSpec(Name string) (string, error) { // return string: json format
	cblogger := cblog.GetLogger("CB-SPIDER")
	cblogger.Info("Mock Driver: called GetOrgVMSpec()!")
	return "", nil
}
