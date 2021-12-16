package mcir

import (
	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
)

type TbSubnetReq struct {
	Description   string                `json:"description"`
	Ipv4_CIDR     string                `json:"ipv4_CIDR"`
	KeyValueInfos []tbcommon.TbKeyValue `json:"keyValueList"`

	Name string `json:"name"`
}

type TbSubnetReqs []TbSubnetReq
