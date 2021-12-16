package mcis

import (
	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
)

type TbVmRecommendInfo struct {
	PlacementAlgo  string                `json:"placementAlgo"`
	PlacementParam []tbcommon.TbKeyValue `json:"placementParam"`
	VmPriority     []TbVmPriority        `json:"vmPriority"`

	VmReq TbVmRecommendReq `json:"vmReq"`
}
