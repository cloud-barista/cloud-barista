package mcis

type MonResultSimple struct {
	Err    string `json:"err"`
	Metric string `json:"metric"` // yes, no
	Value  string `json:"value"`
	VmId   string `json:"vmId"`
}
