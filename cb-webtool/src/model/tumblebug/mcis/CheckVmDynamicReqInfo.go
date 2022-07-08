package mcis

import (
	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
	tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"
)

type CheckVmDynamicReqInfo struct {
	ConnectionConfigCandidates   []string            `json:"connectionConfigCandidates"`
	Region                       tbcommon.TbRegion   `json:"region"`
	SystemMessage                string              `json:"systemMessage"`
	VmSpec                       tbmcir.TbSpecInfo    `json:"vmSpec"`
}