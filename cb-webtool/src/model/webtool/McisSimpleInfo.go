package webtool

// MCIS의 일부정보만 추려서
type McisSimpleInfo struct {
	// ID     string `json:"id"`
	// Name   string `json:"name"`
	// Status string `json:"status"`
	// VMNum  string `json:"vm_num"`

	// mcis.ID, mcis.status, mcis.name, mcis.description
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	McisStatus  string `json:"mcisStatus"`
	Description string `json:"description"`

	InstallMonAgent string `json:"installMonAgent"`
	Label           string `json:"label"`

	ConnectionCount int `json:"connectionCount"`

	// vm_cnt : 해당 mcis의 vm cnt
	// vm_run_cnt, vm_stop_cnt
	VmCount          int            `json:"vmCount"`
	VmStatusNames    string         `json:"vmStatusNames"`
	VmSimpleList     []VmSimpleInfo `json:"vmSimpleList"`
	VmStatusCountMap map[string]int `json:"vmStatusCountMap"`
	// VmRunningCount    int `json:"vnRunningCount"`
	// VmStoppedCount    int `json:"vmStopped"`
	// VmTerminatedCount int `json:"vmTerminated"`

	// csp : 해당 mcis의 connection cnt
	ConnectionConfigProviderMap   map[string]int `json:"connectionConfigProviderMap"`
	ConnectionConfigProviderNames string         `json:"connectionConfigProviderNames"` // 해당 MCIS 등록된 connection의 provider 목록
	// ConnectionConfigProviderNames []string       `json:"connectionConfigProviderNames"` // 해당 MCIS 등록된 connection의 provider 목록
	ConnectionConfigProviderCount int `json:"connectionConfigProviderCount"`
}
type McisSimpleInfos []McisSimpleInfo
