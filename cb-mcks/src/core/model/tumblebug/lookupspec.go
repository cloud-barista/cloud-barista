package tumblebug

import (
	"net/http"
)

type LookupSpec struct {
	Model
	ConnectionName string `json:"connectionName"`
	CspSpecName    string `json:"cspSpecName"`
	SpiderSpecInfo
}

type SpiderSpecInfo struct { // Spider
	// https://github.com/cloud-barista/cb-spider/blob/master/cloud-control-manager/cloud-driver/interfaces/resources/VMSpecHandler.go

	Region string
	Name   string
	VCpu   SpiderVCpuInfo
	Mem    string
	// Gpu    []SpiderGpuInfo

	// KeyValueList []KeyValue
}

type SpiderVCpuInfo struct { // Spider
	Count string
	Clock string // GHz
}

// type SpiderGpuInfo struct { // Spider
// 	Count string
// 	Mfr   string
// 	Model string
// 	Mem   string
// }

func NewLookupSpec(conf string, specName string) *LookupSpec {
	return &LookupSpec{
		ConnectionName: conf,
		CspSpecName:    specName,
	}
}

func (spec *LookupSpec) LookupSpec() error {
	_, err := spec.execute(http.MethodPost, "/lookupSpec", spec, &spec)
	if err != nil {
		return err
	}

	return nil
}
