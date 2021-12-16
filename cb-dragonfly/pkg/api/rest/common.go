package rest

import (
	"github.com/labstack/echo/v4"
)

func SetMessage(msg string) echo.Map {
	responseMsg := echo.Map{}
	responseMsg["message"] = msg
	return responseMsg
}

type SimpleMsg struct {
	Message string `json:"message" example:"Any message"`
}

type AgentType struct {
	NsId        string `json:"ns_id"`
	McisId      string `json:"mcis_id"`
	VmId        string `json:"vm_id"`
	PublicIp    string `json:"public_ip"`
	UserName    string `json:"user_name"`
	SshKey      string `json:"ssh_key"`
	CspType     string `json:"cspType"`
	Port        string `json:"port"`
	ServiceType string `json:"service_type"`
}

type VMOnDemandMetricTags struct {
	McisId string `json:"mcis_id"`
	NsId   string `json:"ns_id"`
	VmId   string `json:"vm_id"`
}

type VMOnDemandMetricValues struct {
	BytesIn  string `json:"bytes_in"`
	BytesOut string `json:"bytes_out"`
	PktsIn   string `json:"pkts_in"`
	pktsOut  string `json:"pkts_out"`
}

type VMOnDemandMetricType struct {
	Name   string                   `json:"name"`
	Tags   []VMOnDemandMetricTags   `json:"tags"`
	Time   string                   `json:"time"`
	Values []VMOnDemandMetricValues `json:"values"`
}

type VMMonInfoTags struct {
	HostId string `json:"host_id"`
}

type VMMonInfoValues struct {
	Free        int     `json:"free"`
	ReadBytes   float64 `json:"read_bytes"`
	ReadTime    float64 `json:"read_time"`
	Reads       float64 `json:"reads"`
	Time        string  `json:"time"`
	Total       int     `json:"total"`
	Used        int     `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	WriteBytes  float64 `json:"write_bytes"`
	WriteTime   float64 `json:"write_time"`
	Writes      float64 `json:"writes"`
}

type VMMonInfoType struct {
	Name   string            `json:"name"`
	Tags   VMMonInfoTags     `json:"tags"`
	Time   string            `json:"time"`
	Values []VMMonInfoValues `json:"values"`
}

// JSONResult's data field will be overridden by the specific type
type JSONResult struct {
	//Code    int          `json:"code" `
	//Message string       `json:"message"`
	//Data    interface{}  `json:"data"`
}
